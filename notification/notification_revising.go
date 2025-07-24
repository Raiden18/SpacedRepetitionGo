package notification

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"strconv"
)

type RevisingNotification struct {
	Boxes      []box.Box
	Flashcards []flashcard.Flashcard
}

func NewEmptyRevisingNotification() RevisingNotification {
	return NewRevisingNotification(
		[]box.Box{},
		[]flashcard.Flashcard{},
	)
}

func NewRevisingNotification(boxes []box.Box, flashcards []flashcard.Flashcard) RevisingNotification {
	return RevisingNotification{
		Boxes:      boxes,
		Flashcards: flashcards,
	}
}

func (revising RevisingNotification) WellDoneMessage() string {
	return "Good Job! 😎 Everything is revised! ✅"
}

func (revising RevisingNotification) TextBody() string {
	return `You have ` + strconv.Itoa(len(revising.Flashcards)) + ` flashcards to revise 🧠`
}

func (revising RevisingNotification) BuildCallback(box box.Box) string {
	return KeyValue(BoxIdToRevise(), box.Id)
}

func (revising RevisingNotification) GetBoxes() []box.Box {
	return revising.Boxes
}

func (revising RevisingNotification) GetFlashCards() []flashcard.Flashcard {
	return revising.Flashcards
}

func (revising RevisingNotification) GetDBTableName() string {
	return "sent_need_revising_notifications"
}
