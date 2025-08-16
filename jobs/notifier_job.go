package jobs

import (
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/telegram"
	"spacedrepetitiongo/utils"
)

func Notifiy() {
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
