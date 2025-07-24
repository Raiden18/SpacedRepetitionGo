package main

import (
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func forgetFlashCardAndShowAnother(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	updateRevisedFlashCardAndShowAnother(
		db,
		update,
		bot,
		client,
		func(flashcardToForget flashcard.Flashcard) flashcard.Flashcard {
			forgottenFlashCard := flashcardToForget.Forget()
			flashcard.InsertFlashCardsIntoDB(db, []flashcard.Flashcard{flashcardToForget}, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
			return forgottenFlashCard
		},
	)

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedMessage(bot, db)
}

func recallFlashCardAndShowAnother(updaet tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	updateRevisedFlashCardAndShowAnother(
		db,
		updaet,
		bot,
		client,
		flashcard.Flashcard.Recall,
	)
}

func updateRevisedFlashCardAndShowAnother(db sqlx.DB, update tgbotapi.Update, bot telegram.Bot, client notion.Client, revisedAction func(flashcard.Flashcard) flashcard.Flashcard) {
	flashCardId := fetchValue(update.CallbackData())
	flashCardToUpdate := flashcard.GetFromDdById(db, flashCardId, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	updatedFlashCard := revisedAction(flashCardToUpdate)
	updatedFlashCard.RemoveFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	bot.DeleteMessage(update.CallbackQuery.Message.MessageID)
	go func() { updatedFlashCard.UpdatePageOnNotion(client) }()
	sendNewFlashCardToTelegramIfExistsToRevise(db, updatedFlashCard.BoxId, bot)

	notification.
		NewRevisingNotificationFromDB(db).
		EditExistedMessage(bot, db)
}

func sendNewFlashCardToTelegramIfExistsToRevise(db sqlx.DB, boxId string, bot telegram.Bot) {
	sendNewFlashCardToTelegramIfExists(
		db,
		boxId,
		flashcard.FLASH_CARDS_TO_REVISE_TABLE,
		func(f flashcard.Flashcard) {
			flashcard.SendToTelegram(bot, flashcard.NewFlashcardTelegramMessageToRevise(f))
		},
	)
}
