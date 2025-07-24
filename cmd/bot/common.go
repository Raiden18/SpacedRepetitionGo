package main

import (
	"spacedrepetitiongo/flashcard"

	"strings"

	"github.com/jmoiron/sqlx"
)

func sendNewFlashCardToTelegramIfExists(db sqlx.DB, boxId string, tableName string, sendToTelegram func(flashcard.Flashcard)) {
	if flashcard.CountInDb(db, boxId, tableName) > 0 {
		flashCard := flashcard.GetFormBox(db, boxId, tableName)
		sendToTelegram(flashCard)
	}
}

func fetchValue(payload string) string {
	return strings.Split(payload, "=")[1]
}
