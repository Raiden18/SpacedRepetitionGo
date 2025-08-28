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
	toPreviosAndNext := []tgbotapi.InlineKeyboardButton{}
	toStartAndEnd := []tgbotapi.InlineKeyboardButton{}
	if message.Flashcard.Previous != nil {
		toPreviosAndNext = append(toPreviosAndNext, previousButton(message.Flashcard))
		toStartAndEnd = append(toStartAndEnd, toBeginningButton(message.Flashcard))
	}
	if message.Flashcard.Next != nil {
		toPreviosAndNext = append(toPreviosAndNext, newNextButton(message.Flashcard))
		toStartAndEnd = append(toStartAndEnd, toEndButton(message.Flashcard))
	}
	rows := [][]tgbotapi.InlineKeyboardButton{
		toPreviosAndNext,
		toStartAndEnd,
		tgbotapi.NewInlineKeyboardRow(
			newMemorizedButton(message.Flashcard),
			finishMemorizingButton(message.Flashcard),
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

func previousButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"⬅️ Previous",
		Parameter(
			PreviousMemorizingFlashCardKey(),
			*flashcard.Previous,
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

func finishMemorizingButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Finish 🏁",
		Parameter(
			FinishMemorizinKey(),
			flashcard.Id,
		),
	)
}

func toEndButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"To the end ⏭️",
		Parameter(
			EndKey(),
			flashcard.Id,
		),
	)
}

func toBeginningButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"To the beginning ⏮️",
		Parameter(
			BeginingKey(),
			flashcard.Id,
		),
	)
}

func NextMemorizingFlashCardKey() string {
	return "nextFlashCardId"
}

func PreviousMemorizingFlashCardKey() string {
	return "previousFlashCardId"
}

func MemorizedMemorizingFlashCardKey() string {
	return "memorizeFlashCardId"
}

func FinishMemorizinKey() string {
	return "finishMemorizing"
}

func EndKey() string {
	return "endFlashcard"
}

func BeginingKey() string {
	return "beginingFlashcard"
}
