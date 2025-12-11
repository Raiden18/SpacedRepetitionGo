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

func (message *FlashcardTelegramMessageToRevise) SendToTelegram(bot telegram.Bot) {
	if message == nil {
		return
	}
	SendToTelegram(bot, message)
}

func (message FlashcardTelegramMessageToRevise) GetButtons() *tgbotapi.InlineKeyboardMarkup {
	toPreviosAndNext := []tgbotapi.InlineKeyboardButton{}
	toStartAndEnd := []tgbotapi.InlineKeyboardButton{}
	if message.Flashcard.Previous != nil {
		toPreviosAndNext = append(toPreviosAndNext, previousButton(message.Flashcard))
		toStartAndEnd = append(toStartAndEnd, toBeginningButton(message.Flashcard))
	}
	if message.Flashcard.Next != nil {
		toPreviosAndNext = append(toPreviosAndNext, nextButton(message.Flashcard))
		toStartAndEnd = append(toStartAndEnd, toEndButton(message.Flashcard))
	}

	rows := [][]tgbotapi.InlineKeyboardButton{
		toPreviosAndNext,
		toStartAndEnd,
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
