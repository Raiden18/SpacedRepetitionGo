package notification

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/telegram"
	"strconv"

	"github.com/jmoiron/sqlx"
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

func NewMemorizingNotificationFromDB(db sqlx.DB) MemorizingNotification {
	return NewMemorizingNotification(
		box.NewBoxesFromDb(db),
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
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

func (memorizing MemorizingNotification) EditExistedMessage(bot telegram.Bot, db sqlx.DB) {
	EditExistedMessage(bot, db, memorizing)
}

func (memorizing MemorizingNotification) EditExistedNotificationOrSendNewIfNotSent(db sqlx.DB, bot telegram.Bot) {
	previouslySentMessageId := GetSentMessageId(db, memorizing.GetDBTableName())
	if previouslySentMessageId == nil {
		SendNewNotification(bot, db, memorizing)
	} else {
		EditExistedMessage(bot, db, memorizing)
	}
}
