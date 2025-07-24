package notification

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
)

type Notification interface {
	WellDoneMessage() string
	TextBody() string
	BuildCallback(box box.Box) string
	GetBoxes() []box.Box
	GetFlashCards() []flashcard.Flashcard
	GetDBTableName() string
}
