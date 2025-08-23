package flashcard

import (
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FlashcardTelegramMessageToMemorize struct {
	Flashcard Flashcard
}

func NewFlashcardTelegramMessageToMemorize(flashcard Flashcard) FlashcardTelegramMessageToMemorize {
	return FlashcardTelegramMessageToMemorize{
		Flashcard: flashcard,
	}
}

func (message FlashcardTelegramMessageToMemorize) GetButtons() *tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			newMemorizedButton(message.Flashcard),
			newNextButton(message.Flashcard),
		),
		tgbotapi.NewInlineKeyboardRow(
			newStartOverButton(message.Flashcard),
		),
	}
	externalButton := createExtranalButton(message.Flashcard)
	if externalButton != nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*externalButton))
	}

	buttons := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &buttons
}

func (message FlashcardTelegramMessageToMemorize) GetFlashCard() Flashcard {
	return message.Flashcard
}

func (message *FlashcardTelegramMessageToMemorize) SendToTelegram(bot telegram.Bot) {
	if message == nil {
		return
	}
	SendToTelegram(bot, message)
}

func newMemorizedButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Memorized ✅",
		Parameter(
			MemorizedMemorizingFlashCardKey(),
			flashcard.Id,
		),
	)
}

func newNextButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Next ➡️",
		Parameter(
			NextMemorizingFlashCardKey(),
			*flashcard.Next,
		),
	)
}

func newStartOverButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Start over 🔄",
		Parameter(
			StartOverMemorizingFlashCardKey(),
			flashcard.Id,
		),
	)
}

func NextMemorizingFlashCardKey() string {
	return "nextFlashCardId"
}

func MemorizedMemorizingFlashCardKey() string {
	return "memorizeFlashCardId"
}

func StartOverMemorizingFlashCardKey() string {
	return "restartMemFlashCardId"
}
