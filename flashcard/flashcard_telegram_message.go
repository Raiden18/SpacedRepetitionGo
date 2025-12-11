package flashcard

import (
	"spacedrepetitiongo/config"

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
