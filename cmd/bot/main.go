package main

import (
	"log"
	"net/http"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	bot := telegram.NewBot()
	notionClient := notion.NewClient()
	db := utils.OpenDb()

	updates := bot.ListenForWebhook()

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		if update.CallbackQuery != nil {
			respondToPressedButtons(update, bot, db, notionClient)
		}
	}
}

func respondToPressedButtons(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	key := fetchKey(update.CallbackQuery.Data)
	buttons := createButtons()
	pressedButtonCallback := buttons[key]
	pressedButtonCallback(update, bot, db, client)
	bot.ResponseToPressedButton(update.CallbackQuery)
}

type ButtonCallbackFunc func(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client)

func createButtons() map[string]ButtonCallbackFunc {
	return map[string]ButtonCallbackFunc{
		notification.BoxIdToRevise():                showFlashCardForSelectedBoxToRevise,
		notification.BoxIdToMemorize():              showFlashCardForSelectedBoxToMemorize,
		flashcard.ForgottenFlashCardKey():           forgetFlashCardAndShowAnother,
		flashcard.RecalledFlashcardKey():            recallFlashCardAndShowAnother,
		flashcard.NextMemorizingFlashCardKey():      hideCurrentMemorizingFlashCardAndShowNext,
		flashcard.StartOverMemorizingFlashCardKey(): startOvertMemorizing,
		flashcard.MemorizedMemorizingFlashCardKey(): memorizeFlashCardAndShowNext,
	}
}

func fetchKey(payload string) string {
	return strings.Split(payload, "=")[0]
}
