package flashcard

import (
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FlashcardTelegramMessageToRevise struct {
	Flashcard Flashcard
}

func NewFlashcardTelegramMessageToRevise(flashcard Flashcard) FlashcardTelegramMessageToRevise {
	return FlashcardTelegramMessageToRevise{
		Flashcard: flashcard,
	}
}

func (message FlashcardTelegramMessageToRevise) GetButtons() *tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			newForgotButton(message.Flashcard),
			newRecalledButton(message.Flashcard),
		),
	}
	externalButton := createExtranalButton(message.Flashcard)
	if externalButton != nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*externalButton))
	}

	buttons := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &buttons
}

func (message FlashcardTelegramMessageToRevise) GetFlashCard() Flashcard {
	return message.Flashcard
}

func newForgotButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Forgot ❌",
		Parameter(
			ForgottenFlashCardKey(),
			flashcard.Id,
		),
	)
}

func newRecalledButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Recalled ✅",
		Parameter(
			RecalledFlashcardKey(),
			flashcard.Id,
		),
	)
}

func ForgottenFlashCardKey() string {
	return "forgottenFlashCardId"
}

func RecalledFlashcardKey() string {
	return "recalledFlashCardId"
}
