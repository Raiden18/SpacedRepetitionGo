package bot

import (
	"spacedrepetitiongo/jobs"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	"github.com/jmoiron/sqlx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandCallbackFunc func(update *tgbotapi.Message, bot telegram.Bot, db sqlx.DB, client notion.Client)

func CreateCommands() map[string]CommandCallbackFunc {
	return map[string]CommandCallbackFunc{
		"update":             updateDataBases,
		"greek_active_vocab": generateGreekActiveVocabulary,
	}
}

func updateDataBases(message *tgbotapi.Message, bot telegram.Bot, db sqlx.DB, client notion.Client) {
	jobs.Update()
	jobs.ReplaceImages()
	jobs.Notifiy()
	bot.DeleteMessage(
		message.MessageID,
	)
}

func generateGreekActiveVocabulary(message *tgbotapi.Message, tg telegram.Bot, db sqlx.DB, client notion.Client) {
	SaveState(db, NewAddGreekVocabularyState())
}
