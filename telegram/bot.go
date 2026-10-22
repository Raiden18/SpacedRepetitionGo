package telegram

import (
	"crypto/subtle"
	"log"
	"net/http"
	"spacedrepetitiongo/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	BLUE_SCREEN = "https://neosmart.net/wiki/wp-content/uploads/sites/5/2013/08/unmountable-boot-volume.png"
)

type Bot struct {
	ImplBot *tgbotapi.BotAPI
}

func NewBot() Bot {
	impl, error := tgbotapi.NewBotAPI(
		config.TelegramApiKey(),
	)
	if error != nil {
		log.Fatal("Could not create Telegram Bot: ", error)
	}
	return Bot{
		ImplBot: impl,
	}
}

func (bot Bot) SendTextMessage(
	text string,
	buttons *tgbotapi.InlineKeyboardMarkup,
) tgbotapi.Message {
	message, error := bot.ImplBot.Send(
		NewMessageConfig(
			text,
			buttons,
		),
	)

	if error != nil {
		log.Println("Could not send new notification to Chat.", error)
	}

	return message
}

func (bot Bot) SendPhotoMessage(
	text string,
	photo string,
	buttons *tgbotapi.InlineKeyboardMarkup,
) {
	photoConfig := PhotomessageConfig(
		text,
		photo,
		buttons,
	)
	_, error := bot.ImplBot.Send(photoConfig)
	if error != nil {
		log.Println("Could not send image: "+photo, error)
		photoConfig.File = tgbotapi.FileURL(BLUE_SCREEN)
		_, error := bot.ImplBot.Send(photoConfig)
		if error != nil {
			log.Println("Could not send image event with BLUE SCREEN. Photo:"+photo, error)
		}
	}
}

func (bot Bot) EditTextMessage(
	messageId int,
	text string,
	buttons *tgbotapi.InlineKeyboardMarkup,
) {
	_, error := bot.ImplBot.Send(
		NewEditMessageTextConfig(
			messageId,
			text,
			buttons,
		),
	)

	if error != nil {
		log.Println("Could edit old notification in chat.", error)
	}
}

// ListenForWebhook serves Telegram updates on /webhook. When secret is set,
// requests without Telegram's matching X-Telegram-Bot-Api-Secret-Token header
// are rejected.
func (bot Bot) ListenForWebhook(secret string) tgbotapi.UpdatesChannel {
	updates := make(chan tgbotapi.Update, bot.ImplBot.Buffer)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if secret != "" && subtle.ConstantTimeCompare(
			[]byte(r.Header.Get("X-Telegram-Bot-Api-Secret-Token")),
			[]byte(secret),
		) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		update, err := bot.ImplBot.HandleUpdate(r)
		if err != nil {
			log.Println("Could not parse webhook update:", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		updates <- *update
	})

	return updates
}

func (bot Bot) ResponseToPressedButton(callback *tgbotapi.CallbackQuery) {
	response := tgbotapi.NewCallback(callback.ID, "")
	_, error := bot.ImplBot.Request(response)
	if error != nil {
		log.Println("Could not Response to pressed button.", error)
	}
}

func (bot Bot) DeleteMessage(messageId int) {
	_, error := bot.ImplBot.Request(
		NewDeleteMessageConfig(messageId),
	)
	if error != nil {
		log.Println("Could not delete message from chat.", error)
	}
}
