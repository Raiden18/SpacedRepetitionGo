package main

import (
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := utils.OpenDb()
	bot := telegram.NewBot()

	notification.
		NewRevisingNotificationFromDB(db).
		SendNewNotificationAndDeleteOld(db, bot)

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedNotificationOrSendNewIfNotSent(db, bot)

	defer db.Close()
}
