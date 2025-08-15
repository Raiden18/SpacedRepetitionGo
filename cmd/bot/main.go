package main

import (
	"log"
	"net/http"
	"spacedrepetitiongo/bot"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/openai"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	tg := telegram.NewBot()
	notionClient := notion.NewClient()
	db := utils.OpenDb()
	openAi := openai.NewClient()

	updates := tg.ListenForWebhook()
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		if update.CallbackQuery != nil {
			bot.RespondToPressedButtons(update, tg, db, notionClient)
		}
		message := update.Message
		if message != nil {
			if message.IsCommand() {
				command := update.Message.Command()
				commands := bot.CreateCommands()
				commands[command](message, tg, db, notionClient)
			}
			if bot.IsAddGreekVocabularyState(db) {
				gptFlashCard := flashcard.GenerateFromGPT(
					"Greek",
					"αρχίζω",
					openAi,
				)
				toMemorize := flashcard.NewFlashcardTelegramMessageToMemorize(gptFlashCard)
				toMemorize.SendToTelegram(tg)
			}
		}
	}
}
