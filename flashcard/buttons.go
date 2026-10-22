package flashcard

import (
	"net/url"
	"spacedrepetitiongo/telegram"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func newCambridgeDictionaryButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewUrlButton(
		"Cambridge Dictionary ⬆️",
		"https://dictionary.cambridge.org/dictionary/english/"+url.PathEscape(flashcard.Name),
	)
}

func newForvoButton(text string) tgbotapi.InlineKeyboardButton {
	replacer := strings.NewReplacer(
		"ο ", "",
		"η ", "",
		"το ", "",
		"\n", "",
		"||", "",
		"ο/η", "",
		"η/ο", "",
	)
	url := "https://forvo.com/search/" + url.PathEscape(replacer.Replace(text))
	return telegram.NewUrlButton(
		"Forvo ⬆️",
		url,
	)
}

func nextButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Next ➡️",
		Parameter(
			NextFlashCardKey(),
			*flashcard.Next,
		),
	)
}

func previousButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"⬅️ Previous",
		Parameter(
			PreviousFlashCardKey(),
			*flashcard.Previous,
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

func finishButton(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Finish 🏁",
		Parameter(
			FinishKey(),
			flashcard.Id,
		),
	)
}

func Parameter(key string, value string) string {
	return key + "=" + value
}

func NextFlashCardKey() string {
	return "nextFlashCardId"
}

func PreviousFlashCardKey() string {
	return "previousFlashCardId"
}

func FinishKey() string {
	return "finish"
}

func EndKey() string {
	return "toEnd"
}

func BeginingKey() string {
	return "toBegining"
}
