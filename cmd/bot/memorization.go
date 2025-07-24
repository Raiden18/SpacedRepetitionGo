package main

import (
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func showFlashCardForSelectedBoxToMemorize(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
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

	flashcard.
		NewRevisingFlashcardFromDbByBoxId(db, memorizedFlashCard.BoxId).
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
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

	flashcard.
		NewMemorizingFlashcardFromDbByBoxId(db, selectedFlashCard.BoxId).
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}

func resetMemorizingProcess(db sqlx.DB, boxId string, bot telegram.Bot) {
	flashCards := flashcard.GetAllFromBdByBoxId(db, boxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.ClearFlashCardTable(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	flashcard.InsertFlashCardsIntoDB(db, flashCards, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)

	flashcard.
		NewMemorizingFlashcardFromDbByBoxId(db, boxId).
		ToTelegramMessageToMemorize().
		SendToTelegram(bot)
}
