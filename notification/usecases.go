package notification

import (
	"spacedrepetitiongo/telegram"

	"github.com/jmoiron/sqlx"
)

func DeleteOldAndSendNewMessageOrEditToDone(
	bot telegram.Bot,
	notification Notification,
	db sqlx.DB,
	messageId int,
) {
	if len(notification.GetFlashCards()) > 0 {
		DeleteOldNotification(bot, db, notification, messageId)
		SendNewNotification(bot, db, notification)
	} else {
		EditNotificationToWellDone(bot, notification, messageId)
	}
}

func EditExistedMessage(
	bot telegram.Bot,
	db sqlx.DB,
	notification Notification,
) {
	messageId := *GetSentMessageId(db, notification.GetDBTableName())
	if len(notification.GetFlashCards()) > 0 {
		EditNotificationWithButtons(bot, notification, messageId)
	} else {
		EditNotificationToWellDone(bot, notification, messageId)
	}
}

func SendNewNotification(
	bot telegram.Bot,
	db sqlx.DB,
	notification Notification,
) {
	message := SendToChat(bot, notification)
	InsertIntoDb(db, message.MessageID, notification.GetDBTableName())
}

func DeleteOldNotification(
	bot telegram.Bot,
	db sqlx.DB,
	notification Notification,
	messageId int,
) {
	bot.DeleteMessage(messageId)
	DeleteFromDb(db, messageId, notification.GetDBTableName())
}
