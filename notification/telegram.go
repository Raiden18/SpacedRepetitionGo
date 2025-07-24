package notification

import (
	"spacedrepetitiongo/telegram"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendToChat(bot telegram.Bot, notification Notification) tgbotapi.Message {
	return bot.SendTextMessage(
		notification.TextBody(),
		buttons(notification),
	)
}

func EditNotificationToWellDone(bot telegram.Bot, notification Notification, messageId int) {
	bot.EditTextMessage(
		messageId,
		notification.WellDoneMessage(),
		nil,
	)
}

func EdittNotificationWithButtons(bot telegram.Bot, notification Notification, messageId int) {
	bot.EditTextMessage(
		messageId,
		notification.TextBody(),
		buttons(notification),
	)
}

func buttons(notification Notification) *tgbotapi.InlineKeyboardMarkup {
	var replyMarkups [][]tgbotapi.InlineKeyboardButton

	for _, aBox := range notification.GetBoxes() {
		flashcardsFromBox := aBox.GetFlashCardsFromBox(notification.GetFlashCards())
		amount := len(flashcardsFromBox)
		callback := notification.BuildCallback(aBox)
		if amount > 0 {
			boxButton := boxButton(aBox.Name, amount, callback)
			replyMarkups = append(
				replyMarkups,
				tgbotapi.NewInlineKeyboardRow(boxButton),
			)
		}
	}

	buttons := tgbotapi.NewInlineKeyboardMarkup(replyMarkups...)
	return &buttons
}

func boxButton(name string, count int, callback string) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		name+": "+strconv.Itoa(count),
		callback,
	)
}
