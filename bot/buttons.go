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
		flashcard.PreviousMemorizingFlashCardKey():  onPreviousButtonOfMemorizingFlashcardClicked,
		flashcard.FinishMemorizinKey():              onFinishedMemorizingButtonClicked,
		flashcard.MemorizedMemorizingFlashCardKey(): onMemorizedButtonOfFlashcardClicked,
	}
}

func onBoxButtonToReviseClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	sendNextFlashcardToRevise(db, bot, fetchValue(update.CallbackData()))
}

func onForgetButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	flashcardToForget := flashcard.NewRevisingFlashcardcFromDbById(db, fetchValue(update.CallbackData()))
	forgottenFlashcard := flashcardToForget.Forget()

	if flashcard.CountInDb(db, forgottenFlashcard.BoxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE) > 0 {
		last := flashcard.GetLast(db, forgottenFlashcard.BoxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
		last.Next = &forgottenFlashcard.Id
		last.UpdateLinkedFlashCards(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
		forgottenFlashcard.Previous = &last.Id
	} else {
		forgottenFlashcard.Previous = nil
	}
	forgottenFlashcard.Next = nil

	forgottenFlashcard.
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
	boxId := fetchValue(update.CallbackData())
	flashCards := flashcard.GetAllFromBdByBoxId(db, boxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.ClearTable(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	flashcard.InsertIntoDB(db, flashCards, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	firstFlashCard := flashcard.GetFirst(db, boxId, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	firstFlashCard.
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}

func onMemorizedButtonOfFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, notionClient notion.Client) {
	memorizedFlashCard := flashcard.
		NewMemorizingFlashcardFromDb(db, fetchValue(update.CallbackData())).
		Memorize().
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE).
		RemoveFrom(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE).
		UpdateOnNotion(notionClient).
		RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)

	if memorizedFlashCard.Previous != nil {
		previous := flashcard.NewMemorizingFlashcardFromDb(db, *memorizedFlashCard.Previous)
		previous.Next = memorizedFlashCard.Next
		previous.UpdateLinkedFlashCards(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
		previous.UpdateLinkedFlashCards(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	}

	if memorizedFlashCard.Next != nil {
		next := flashcard.NewMemorizingFlashcardFromDb(db, *memorizedFlashCard.Next)
		next.Previous = memorizedFlashCard.Previous
		next.UpdateLinkedFlashCards(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
		next.UpdateLinkedFlashCards(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	}

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedMessage(db, bot)

	if memorizedFlashCard.Next != nil {
		flashcard.
			NewMemorizingFlashcardFromDb(db, *memorizedFlashCard.Next).
			ToTelegramMessageToMemorize().
			SendToTelegram(bot)
	} else if memorizedFlashCard.Previous != nil {
		flashcard.
			NewMemorizingFlashcardFromDb(db, *memorizedFlashCard.Previous).
			ToTelegramMessageToMemorize().
			SendToTelegram(bot)
	}
}

func onFinishedMemorizingButtonClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, notionClient notion.Client) {
	bot.DeleteMessage(update.CallbackQuery.Message.MessageID)
}

func onPreviousButtonOfMemorizingFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	previousFlashcardId := fetchValue(update.CallbackData())
	previousFlashCard := flashcard.NewMemorizingFlashcardFromDb(db, previousFlashcardId)

	currentFlashcard := flashcard.NewMemorizingFlashcardFromDb(db, *previousFlashCard.Next)
	currentFlashcard.RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)

	previousFlashCard.
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}

func onNextButtonOfMemorizingFlashcardClicked(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	nextFlashcardId := fetchValue(update.CallbackData())
	nextFlashCard := flashcard.NewMemorizingFlashcardFromDb(db, nextFlashcardId)

	currentFlashcard := flashcard.NewMemorizingFlashcardFromDb(db, *nextFlashCard.Previous)
	currentFlashcard.RemoveFromChat(bot, update.CallbackQuery.Message.MessageID)

	nextFlashCard.
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}
