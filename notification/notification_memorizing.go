package notification

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"strconv"
)

type MemorizingNotification struct {
	Boxes      []box.Box
	Flashcards []flashcard.Flashcard
}

func NewEmptyMemorizingNotification() MemorizingNotification {
	return NewMemorizingNotification(
		[]box.Box{},
		[]flashcard.Flashcard{},
	)
}

func NewMemorizingNotification(
	boxes []box.Box,
	flashcards []flashcard.Flashcard,
) MemorizingNotification {
	return MemorizingNotification{
		Boxes:      boxes,
		Flashcards: flashcards,
	}
}

func (memorizing MemorizingNotification) WellDoneMessage() string {
	return "Good Job! 😎 Everything is memorized! ✅"
}

func (memorizing MemorizingNotification) TextBody() string {
	return `You have ` + strconv.Itoa(len(memorizing.Flashcards)) + ` flashcards to memorize 🧠`
}

func (memorizing MemorizingNotification) BuildCallback(box box.Box) string {
	return KeyValue(BoxIdToMemorize(), box.Id)
}

func (memorizing MemorizingNotification) GetBoxes() []box.Box {
	return memorizing.Boxes
}

func (memorizing MemorizingNotification) GetFlashCards() []flashcard.Flashcard {
	return memorizing.Flashcards
}

func (memorizing MemorizingNotification) GetDBTableName() string {
	return "sent_need_memorizing_notification"
}
