package main

import (
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func showFlashCardForSelectedBoxToMemorize(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedBoxId := fetchValue(update.CallbackData())
	resetMemorizingProcess(db, selectedBoxId, bot)
}

func showFlashCardForSelectedBoxToRevise(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	selectedBoxId := fetchValue(update.CallbackQuery.Data)
	sendNewFlashCardToTelegramIfExistsToRevise(db, selectedBoxId, bot)
}
