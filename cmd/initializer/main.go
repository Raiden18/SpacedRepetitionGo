package main

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := utils.OpenDb()
	box.CreateTableIfNotExist(db)

	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	flashcard.CreateTableIfNotExist(db, flashcard.FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)

	for _, emptyNotification := range createNotifications() {
		notification.CreateTableIfNotExist(db, emptyNotification.GetDBTableName())
	}

	db.Close()
}

func createNotifications() []notification.Notification {
	return []notification.Notification{
		notification.NewEmptyRevisingNotification(),
		notification.NewEmptyRevisingNotification(),
	}
}
