package telegram

import (
	"log"
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
		log.Println("Could not create Telegram Bot", error)
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
		log.Fatalln("Could not send new notification to Chat.", error)
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
			log.Fatalln("Could not send image event with BLUE SCREEN.", error)
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

func (bot Bot) ListenForWebhook() tgbotapi.UpdatesChannel {
	return bot.ImplBot.ListenForWebhook("/webhook")
}

func (bot Bot) ResponseToPressedButton(callback *tgbotapi.CallbackQuery) {
	response := tgbotapi.NewCallback(callback.ID, "")
	bot.ImplBot.Request(response)
}

func (bot Bot) DeleteMessage(messageId int) {
	_, error := bot.ImplBot.Request(
		NewDeleteMessageConfig(messageId),
	)
	if error != nil {
		log.Println("Could not delete message from chat.", error)
	}
}
