package notification

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type MessageId struct {
	Id int `db:"message_id"`
}

func GetSentMessageId(db sqlx.DB, tableName string) *int {
	var messages = []*MessageId{}
	error := db.Select(&messages, `SELECT * FROM `+tableName)
	if error != nil {
		log.Fatal(error)
	}
	if len(messages) == 0 {
		return nil
	}
	return &messages[0].Id
}

func InsertIntoDb(db sqlx.DB, id int, tableName string) {
	query := `INSERT INTO ` + tableName + ` (message_id) VALUES (:message_id);`
	messageId := MessageId{
		Id: id,
	}
	_, err := db.NamedExec(query, &messageId)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteFromDb(db sqlx.DB, id int, tableName string) {
	_, err := db.Exec(`DELETE FROM ` + tableName)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTableIfNotExist(db sqlx.DB, tableName string) {
	_, err := db.Exec(`
    	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
    	    message_id BIGINT NOT NULL PRIMARY KEY
    	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
