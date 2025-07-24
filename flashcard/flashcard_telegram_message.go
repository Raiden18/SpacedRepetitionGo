package flashcard

import (
	"net/url"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/telegram"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FlashcardTelegramMessage interface {
	GetButtons() *tgbotapi.InlineKeyboardMarkup
	GetFlashCard() Flashcard
}

func createExtranalButton(flashcard Flashcard) *tgbotapi.InlineKeyboardButton {
	if config.GreekVocabularyId() == flashcard.BoxId {
		forvoButton := newForvoButton(*flashcard.Explanation)
		return &forvoButton
	}
	if config.GreekPassiveVocabularyId() == flashcard.BoxId {
		forvoButton := newForvoButton(flashcard.Name)
		return &forvoButton
	}
	if config.EnglishVocabularyId() == flashcard.BoxId {
		cambridgeButton := newCambridgeDictionaryButton(flashcard)
		return &cambridgeButton
	}
	return nil
}

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

func Parameter(key string, value string) string {
	return key + "=" + value
}
