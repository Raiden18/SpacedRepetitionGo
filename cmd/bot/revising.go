package main

import (
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func showFlashCardForSelectedBoxToRevise(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	flashcard.
		NewRevisingFlashcardFromDbByBoxId(db, fetchValue(update.CallbackData())).
		ToTelegramMessageToRevise().
		SendToTelegram(bot)
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

	flashcard.
		NewRevisingFlashcardFromDbByBoxId(db, forgottenFlashcard.BoxId).
		ToTelegramMessageToRevise().
		SendToTelegram(bot)
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

	flashcard.
		NewRevisingFlashcardFromDbByBoxId(db, recalledFlashcard.BoxId).
		ToTelegramMessageToRevise().
		SendToTelegram(bot)
}
