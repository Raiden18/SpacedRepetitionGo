package box

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var (
	BOX_TABLE_NAME = "notion_data_bases"
)

func NewBoxesFromDb(db sqlx.DB) []Box {
	query := "SELECT * FROM " + BOX_TABLE_NAME
	boxes := []Box{}
	db.Select(&boxes, query)
	return boxes
}

func CreateTableIfNotExist(db sqlx.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + BOX_TABLE_NAME + ` (id VARCHAR(255) NOT NULL PRIMARY KEY,text TEXT);`)
	if err != nil {
		log.Fatal(err)
	}
}

func ClearTable(db sqlx.DB) {
	_, err := db.Exec(`DELETE FROM ` + BOX_TABLE_NAME)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertIntoDB(db sqlx.DB, boxes []Box) {
	for _, box := range boxes {
		insertBox(db, box)
	}
}

func insertBox(db sqlx.DB, box Box) {
	_, err := db.NamedExec(
		`INSERT INTO `+BOX_TABLE_NAME+` (id, text) VALUES (:id, :text);`,
		box,
	)
	if err != nil {
		log.Fatal(err)
	}
}
