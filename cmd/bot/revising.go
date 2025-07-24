package main

import (
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

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
