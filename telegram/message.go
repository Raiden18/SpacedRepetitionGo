package telegram

import (
	"spacedrepetitiongo/config"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMessageConfig(
	text string,
	buttons *tgbotapi.InlineKeyboardMarkup,
) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           config.TelegramChatId(),
			ReplyToMessageID: 0,
			ReplyMarkup:      buttons,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeMarkdownV2,
		DisableWebPagePreview: false,
	}
}

func NewEditMessageTextConfig(
	messageId int,
	newText string,
	newButtons *tgbotapi.InlineKeyboardMarkup,
) tgbotapi.EditMessageTextConfig {
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      config.TelegramChatId(),
			MessageID:   messageId,
			ReplyMarkup: newButtons,
		},
		Text: newText,
	}
}

func NewDeleteMessageConfig(messageId int) tgbotapi.DeleteMessageConfig {
	return tgbotapi.NewDeleteMessage(
		config.TelegramChatId(),
		messageId,
	)
}

func PhotomessageConfig(
	text string,
	image string,
	buttons *tgbotapi.InlineKeyboardMarkup,
) tgbotapi.PhotoConfig {
	return tgbotapi.PhotoConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      config.TelegramChatId(),
				ReplyMarkup: buttons,
			},
			File: photoFile(image),
		},
		Caption:   text,
		ParseMode: tgbotapi.ModeMarkdownV2,
	}
}

// photoFile sends remote images by URL (cards whose image the image job hasn't
// downloaded yet) and everything else as a local file.
func photoFile(image string) tgbotapi.RequestFileData {
	if strings.HasPrefix(image, "https://") || strings.HasPrefix(image, "http://") {
		return tgbotapi.FileURL(image)
	}
	return tgbotapi.FilePath(image)
}
