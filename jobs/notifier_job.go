package jobs

import (
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/schema"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"
)

func Notifiy() {
	db := utils.OpenDb()
	defer db.Close()
	schema.Ensure(db)
	bot := telegram.NewBot()

	notification.
		NewRevisingNotificationFromDB(db).
		SendNewNotificationAndDeleteOld(db, bot)

	notification.
		NewMemorizingNotificationFromDB(db).
		EditExistedNotificationOrSendNewIfNotSent(db, bot)
}
