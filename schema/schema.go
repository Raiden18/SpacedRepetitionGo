// Package schema creates and migrates every MySQL table the project uses.
// It is safe to run on every start, so each binary calls Ensure and a deploy
// never depends on someone remembering to run the initializer.
package schema

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notification"
	"spacedrepetitiongo/outbox"

	"github.com/jmoiron/sqlx"
)

func Ensure(db sqlx.DB) {
	box.CreateTableIfNotExist(db)

	flashcard.CreateTableAndMigrate(db)

	notification.CreateTableIfNotExist(db, notification.NewEmptyRevisingNotification().GetDBTableName())
	notification.CreateTableIfNotExist(db, notification.NewEmptyMemorizingNotification().GetDBTableName())

	outbox.CreateTableIfNotExist(db)
}
