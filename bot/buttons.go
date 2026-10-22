package bot

import (
	"log"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type ButtonCallbackFunc func(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client)

func RespondToPressedButtons(update tgbotapi.Update, tg telegram.Bot, db sqlx.DB, client notion.Client) {
	defer tg.ResponseToPressedButton(update.CallbackQuery)

	key := strings.Split(update.CallbackQuery.Data, "=")[0]
	pressedButtonCallback, ok := createButtons()[key]
	if !ok {
		log.Printf("Ignoring unknown button: %q", update.CallbackQuery.Data)
		return
	}
	pressedButtonCallback(update, tg, db, client)
}

func createButtons() map[string]ButtonCallbackFunc {
	return map[string]ButtonCallbackFunc{
		notification.BoxIdToRevise():      onBoxButtonToReviseClicked,
		notification.BoxIdToMemorize():    onBoxButtonToMemorizeClicked,
		flashcard.ForgottenFlashCardKey(): onForgetButtonOfFlashcardClicked,
		flashcard.RecalledFlashcardKey():  onRecallButtonOfFlashCardClicked,
		flashcard.NextFlashCardKey():      onShowCardButtonClicked,
		flashcard.PreviousFlashCardKey():  onShowCardButtonClicked,
		flashcard.FinishKey():             onFinishButtonClicked,
		flashcard.MemorizedFlashCardKey(): onMemorizedButtonOfFlashcardClicked,
		flashcard.EndKey():                onToEndClicked,
		flashcard.BeginingKey():           onToBeginingClicked,
	}
}

func messageId(update tgbotapi.Update) int {
	return update.CallbackQuery.Message.MessageID
}

func onToEndClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	jumpWithinQueue(update, bot, db, flashcard.LastInBox)
}

func onToBeginingClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	jumpWithinQueue(update, bot, db, flashcard.FirstInBox)
}

// jumpWithinQueue shows the first or last card of the current card's queue.
// Buttons carry the current card's id; buttons on messages sent before the
// cards table existed carry a box id instead, so fall back to that.
func jumpWithinQueue(
	update tgbotapi.Update,
	bot telegram.Bot,
	db sqlx.DB,
	pick func(db sqlx.DB, boxId string, queue string) *flashcard.Flashcard,
) {
	value := fetchValue(update.CallbackData())

	var target *flashcard.Flashcard
	if current := flashcard.LoadById(db, value); current != nil {
		target = pick(db, current.BoxId, current.Queue)
	} else if target = pick(db, value, flashcard.QueueMemorize); target == nil {
		target = pick(db, value, flashcard.QueueRevise)
	}

	target.SendToTelegram(bot)
	bot.DeleteMessage(messageId(update))
}

func onBoxButtonToReviseClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	flashcard.FirstInBox(db, fetchValue(update.CallbackData()), flashcard.QueueRevise).SendToTelegram(bot)
}

func onBoxButtonToMemorizeClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	flashcard.FirstInBox(db, fetchValue(update.CallbackData()), flashcard.QueueMemorize).SendToTelegram(bot)
}

// loadForAnswer returns the card a button refers to, if it is still in the
// expected queue. Otherwise (already answered, e.g. a double tap) it deletes
// the stale message and returns nil.
func loadForAnswer(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, queue string) *flashcard.Flashcard {
	card := flashcard.LoadById(db, fetchValue(update.CallbackData()))
	if card == nil || card.Queue != queue {
		bot.DeleteMessage(messageId(update))
		return nil
	}
	return card
}

func onForgetButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	card := loadForAnswer(update, bot, db, flashcard.QueueRevise)
	if card == nil {
		return
	}

	// Queue the answer for Notion before touching the local cache, so it can't be lost.
	forgotten := card.
		Forget().
		EnqueueNotionUpdate(db).
		MoveToEndOf(db, flashcard.QueueMemorize).
		RemoveFromChat(bot, messageId(update))

	notification.NewRevisingNotificationFromDB(db).EditExistedMessage(db, bot)
	notification.NewMemorizingNotificationFromDB(db).EditExistedMessage(db, bot)

	flashcard.FirstInBox(db, forgotten.BoxId, flashcard.QueueRevise).SendToTelegram(bot)
}

func onRecallButtonOfFlashCardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	card := loadForAnswer(update, bot, db, flashcard.QueueRevise)
	if card == nil {
		return
	}

	recalled := card.
		Recall().
		EnqueueNotionUpdate(db).
		Delete(db).
		RemoveFromChat(bot, messageId(update))

	notification.NewRevisingNotificationFromDB(db).EditExistedMessage(db, bot)

	flashcard.FirstInBox(db, recalled.BoxId, flashcard.QueueRevise).SendToTelegram(bot)
}

func onMemorizedButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	card := loadForAnswer(update, bot, db, flashcard.QueueMemorize)
	if card == nil {
		return
	}

	memorized := card.
		Memorize().
		EnqueueNotionUpdate(db).
		Delete(db).
		RemoveFromChat(bot, messageId(update))

	notification.NewMemorizingNotificationFromDB(db).EditExistedMessage(db, bot)

	// Carry on with the neighbouring card, preferring the next one.
	if memorized.Next != nil {
		flashcard.LoadById(db, *memorized.Next).SendToTelegram(bot)
	} else if memorized.Previous != nil {
		flashcard.LoadById(db, *memorized.Previous).SendToTelegram(bot)
	}
}

func onFinishButtonClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	bot.DeleteMessage(messageId(update))
}

// onShowCardButtonClicked handles Next and Previous: both carry the id of the
// card to show.
func onShowCardButtonClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	flashcard.LoadById(db, fetchValue(update.CallbackData())).SendToTelegram(bot)
	bot.DeleteMessage(messageId(update))
}
