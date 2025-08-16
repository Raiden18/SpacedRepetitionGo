package jobs

import (
	"spacedrepetitiongo/bot"
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/utils"
)

func Initialize() {
	db := utils.OpenDb()
	box.CreateTableIfNotExist(db)

	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
	bot.CreateTableIfNotExist(db)

	for _, emptyNotification := range createNotifications() {
		notification.CreateTableIfNotExist(db, emptyNotification.GetDBTableName())
	}

	db.Close()
}

func createNotifications() []notification.Notification {
	return []notification.Notification{
		notification.NewEmptyRevisingNotification(),
		notification.NewEmptyMemorizingNotification(),
	}
}
