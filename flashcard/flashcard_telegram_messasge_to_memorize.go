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
	navigationButtons := []tgbotapi.InlineKeyboardButton{}
	if message.Flashcard.Previous != nil {
		navigationButtons = append(navigationButtons, previousButton(message.Flashcard))
	}
	if message.Flashcard.Next != nil {
		navigationButtons = append(navigationButtons, newNextButton(message.Flashcard))
	}
	rows := [][]tgbotapi.InlineKeyboardButton{
		navigationButtons,
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
