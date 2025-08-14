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

type ButtonCallbackFunc func(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client)
type CommandCallbackFunc func()

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
		message := update.Message
		if message != nil {
			if message.IsCommand() {
				command := update.Message.Command()
				commands := createCommands()
				bot.DeleteMessage(
					message.MessageID,
				)
				commands[command]()
			}
		}
	}
}

func respondToPressedButtons(update tgbotapi.Update, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	key := strings.Split(update.CallbackQuery.Data, "=")[0]
	buttons := createButtons()
	pressedButtonCallback := buttons[key]
	pressedButtonCallback(update, bot, db, client)
	bot.ResponseToPressedButton(update.CallbackQuery)
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

func createCommands() map[string]CommandCallbackFunc {
	return map[string]CommandCallbackFunc{
		"update": updateCommand,
	}
}
