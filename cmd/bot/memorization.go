package main

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func memorizeFlashCardAndShowNext(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedFlashCardId := fetchValue(update.CallbackData())
	selectedFlashCard := flashcard.GetFromDdById(db, selectedFlashCardId, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	memorizedFlashCard := selectedFlashCard.Memorize()
	memorizedFlashCard.RemoveFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	memorizedFlashCard.RemoveFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	bot.DeleteMessage(update.CallbackQuery.Message.MessageID)
	go func() { memorizedFlashCard.UpdatePageOnNotion(client) }()
	sendNewFlashCardToTelegramIfExistsToMemorize(db, selectedFlashCard.BoxId, bot)

	memorizingNotification := notification.NewMemorizingNotification(
		box.NewBoxesFromDb(db),
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
	)

	notification.EditExistedMessage(
		bot,
		memorizingNotification,
		*notification.GetSentMessageId(db, memorizingNotification.GetDBTableName()),
	)
}

func startOvertMemorizing(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedFlashCardId := fetchValue(update.CallbackData())
	selectedFlashCard := flashcard.GetFromDdById(db, selectedFlashCardId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	resetMemorizingProcess(db, selectedFlashCard.BoxId, bot)
	bot.DeleteMessage(update.CallbackQuery.Message.MessageID)
}

func hideCurrentMemorizingFlashCardAndShowNext(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedFlashCardId := fetchValue(update.CallbackData())
	selectedFlashCard := flashcard.GetFromDdById(db, selectedFlashCardId, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	selectedFlashCard.RemoveFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	bot.DeleteMessage(update.CallbackQuery.Message.MessageID)
	sendNewFlashCardToTelegramIfExistsToMemorize(db, selectedFlashCard.BoxId, bot)
}

func resetMemorizingProcess(db sqlx.DB, boxId string, tgBot telegram.Bot) {
	flashCards := flashcard.GetAllFromBdByBoxId(db, boxId, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.ClearFlashCardTable(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	flashcard.InsertFlashCardsIntoDB(db, flashCards, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	sendNewFlashCardToTelegramIfExistsToMemorize(db, boxId, tgBot)
}

func sendNewFlashCardToTelegramIfExistsToMemorize(db sqlx.DB, boxId string, tgBot telegram.Bot) {
	sendNewFlashCardToTelegramIfExists(
		db,
		boxId,
		flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE,
		func(f flashcard.Flashcard) {
			flashcard.SendToTelegram(tgBot, flashcard.NewFlashcardTelegramMessageToMemorize(f))
		},
	)
}
