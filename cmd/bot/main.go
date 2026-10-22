package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"spacedrepetitiongo/bot"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/outbox"
	"spacedrepetitiongo/schema"

	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	tg := telegram.NewBot()
	notionClient := notion.NewClient()
	db := utils.OpenDb()
	schema.Ensure(db)
	ownerChatId := config.TelegramChatId()

	outbox.StartWorker(db, notionClient)

	webhookSecret := config.TelegramWebhookSecret()
	if webhookSecret == "" {
		log.Println("WARNING: telegram.webhook_secret is not set; webhook calls are not authenticated")
	}
	updates := tg.ListenForWebhook(webhookSecret)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		handleUpdate(update, ownerChatId, tg, db, notionClient)
	}
}

// handleUpdate processes one update. A failure while handling it is logged
// and the bot carries on with the next update instead of exiting.
func handleUpdate(update tgbotapi.Update, ownerChatId int64, tg telegram.Bot, db sqlx.DB, notionClient notion.Client) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from failure while handling update %d: %v\n%s", update.UpdateID, r, debug.Stack())
		}
	}()

	if update.CallbackQuery != nil {
		message := update.CallbackQuery.Message
		if message == nil || message.Chat == nil || message.Chat.ID != ownerChatId {
			log.Println("Ignoring button press from another chat")
			return
		}
		bot.RespondToPressedButtons(update, tg, db, notionClient)
	}

	message := update.Message
	if message == nil || !message.IsCommand() {
		return
	}
	if message.Chat == nil || message.Chat.ID != ownerChatId {
		log.Printf("Ignoring /%s from another chat", message.Command())
		return
	}
	command := message.Command()
	handler, ok := bot.CreateCommands()[command]
	if !ok {
		log.Printf("Ignoring unknown command: /%s", command)
		return
	}
	handler(message, tg, db, notionClient)
}
