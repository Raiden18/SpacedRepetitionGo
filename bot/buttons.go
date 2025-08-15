package bot

import (
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
	key := strings.Split(update.CallbackQuery.Data, "=")[0]
	buttons := createButtons()
	pressedButtonCallback := buttons[key]
	pressedButtonCallback(update, tg, db, client)
	tg.ResponseToPressedButton(update.CallbackQuery)
}

func createButtons() map[string]ButtonCallbackFunc {
	return map[string]ButtonCallbackFunc{
		notification.BoxIdToRevise():                onBoxButtonToReviseClicked,
		notification.BoxIdToMemorize():              onBoxButtonToMemorizeClicked,
		flashcard.ForgottenFlashCardKey():           onForgetButtonOfFlashcardClicked,
		flashcard.RecalledFlashcardKey():            onRecallButtonOfFlashCardClicked,
		flashcard.NextMemorizingFlashCardKey():      onNextButtonOfMemorizingFlashcardClicked,
		flashcard.StartOverMemorizingFlashCardKey(): onStartOvertButtonOfMemorizingFlashcardClicked,
		flashcard.MemorizedMemorizingFlashCardKey(): onMemorizedButtonOfFlashcardClicked,
	}
}

func onBoxButtonToReviseClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	sendNextFlashcardToRevise(db, bot, fetchValue(update.CallbackData()))
}

func onForgetButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	forgottenFlashcard := flashcard.
		NewRevisingFlashcardcFromDbById(db, fetchValue(update.CallbackData())).
		Forget().
		InsertInto(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE).
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID).
		UpdateOnNotion(client)

	notification.
		NewRevisingNotificationFromDB(db).
		EditExistedMessage(db, bot)

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedMessage(db, bot)

	sendNextFlashcardToRevise(db, bot, forgottenFlashcard.BoxId)
}

func onRecallButtonOfFlashCardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	recalledFlashcard := flashcard.
		NewRevisingFlashcardcFromDbById(db, fetchValue(update.CallbackData())).
		Recall().
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID).
		UpdateOnNotion(client)

	notification.
		NewRevisingNotificationFromDB(db).
		EditExistedMessage(db, bot)

	sendNextFlashcardToRevise(db, bot, recalledFlashcard.BoxId)
}

func sendNextFlashcardToRevise(db sqlx.DB, bot telegram.Bot, boxId string) {
	flashcard.
		NewRevisingFlashcardFromDbByBoxId(db, boxId).
		ToTelegramMessageToRevise().
		SendToTelegram(bot)
}

func onBoxButtonToMemorizeClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedBoxId := fetchValue(update.CallbackData())
	resetMemorizingProcess(db, selectedBoxId, bot)
}

func onMemorizedButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, notionClient notion.Client) {
	memorizedFlashCard := flashcard.
		NewMemorizingFlashcardFromDb(db, fetchValue(update.CallbackData())).
		Memorize().
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE).
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE).
		UpdateOnNotion(notionClient).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedMessage(db, bot)

	sendNewMemorizingFlashcardToChat(db, bot, memorizedFlashCard.BoxId)
}

func onStartOvertButtonOfMemorizingFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, notionClient notion.Client) {
	selectedFlashCard := flashcard.
		NewMemorizingFlashcardFromDb(db, fetchValue(update.CallbackData())).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)

	resetMemorizingProcess(db, selectedFlashCard.BoxId, bot)
}

func onNextButtonOfMemorizingFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedFlashCard := flashcard.
		NewMemorizingFlashcardFromDb(db, fetchValue(update.CallbackData())).
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)
	sendNewMemorizingFlashcardToChat(db, bot, selectedFlashCard.BoxId)
}

func resetMemorizingProcess(db sqlx.DB, boxId string, bot telegram.Bot) {
	flashCards := flashcard.GetAllFromBdByBoxId(db, boxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.ClearFlashCardTable(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	flashcard.InsertFlashCardsIntoDB(db, flashCards, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	sendNewMemorizingFlashcardToChat(db, bot, boxId)
}

func sendNewMemorizingFlashcardToChat(db sqlx.DB, bot telegram.Bot, boxId string) {
	flashcard.
		NewMemorizingFlashcardFromDbByBoxId(db, boxId).
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}
