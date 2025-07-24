package box

import (
	"spacedrepetitiongo/flashcard"
)

type Box struct {
	Id   string `db:"id"`
	Name string `db:"text"`
}

func (box Box) GetFlashCardsFromBox(flashcards []flashcard.Flashcard) []flashcard.Flashcard {
	flashCardsForBox := []flashcard.Flashcard{}
	for _, flashCard := range flashcards {
		if flashCard.BoxId == box.Id {
			flashCardsForBox = append(flashCardsForBox, flashCard)
		}
	}
	return flashCardsForBox
}
