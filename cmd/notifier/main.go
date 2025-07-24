package main

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db := utils.OpenDb()
	bot := telegram.NewBot()
	boxes := box.NewBoxesFromDb(db)
	revisingNotification := notification.NewRevisingNotification(
		boxes,
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE),
	)
	memorizingNotification := notification.NewMemorizingNotification(
		boxes,
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
	)
	sendNewNotificationAndDeleteOld(db, bot, revisingNotification)
	editExistedNotificationOrSendNewIfNotSent(db, bot, memorizingNotification)
	defer db.Close()
}

func sendNewNotificationAndDeleteOld(db sqlx.DB, bot telegram.Bot, notificationStruct notification.Notification) {
	previouslySentMessageId := notification.GetSentMessageId(db, notificationStruct.GetDBTableName())
	if previouslySentMessageId == nil {
		notification.SendNewNotification(bot, db, notificationStruct)
	} else {
		notification.DeleteOldAndSendNewMessageOrEditToDone(bot, notificationStruct, db, *previouslySentMessageId)
	}
}

func editExistedNotificationOrSendNewIfNotSent(db sqlx.DB, bot telegram.Bot, notificationStruct notification.Notification) {
	previouslySentMessageId := notification.GetSentMessageId(db, notificationStruct.GetDBTableName())
	if previouslySentMessageId == nil {
		notification.SendNewNotification(bot, db, notificationStruct)
	} else {
		notification.EditExistedMessage(bot, notificationStruct, *previouslySentMessageId)
	}
}
