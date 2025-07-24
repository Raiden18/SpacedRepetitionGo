package notification

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/telegram"
	"strconv"

	"github.com/jmoiron/sqlx"
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

func NewRevisingNotificationFromDB(db sqlx.DB) RevisingNotification {
	return NewRevisingNotification(
		box.NewBoxesFromDb(db),
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE),
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

func (revising RevisingNotification) EditExistedMessage(bot telegram.Bot, db sqlx.DB) {
	EditExistedMessage(bot, db, revising)
}

func (revising RevisingNotification) SendNewNotificationAndDeleteOld(db sqlx.DB, bot telegram.Bot) {
	previouslySentMessageId := GetSentMessageId(db, revising.GetDBTableName())
	if previouslySentMessageId == nil {
		SendNewNotification(bot, db, revising)
	} else {
		DeleteOldAndSendNewMessageOrEditToDone(bot, revising, db, *previouslySentMessageId)
	}
}
