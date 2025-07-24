package flashcard

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var (
	FLASH_CARDS_TO_REVISE_TABLE              = "flash_cards_to_revise"
	FLASH_CARDS_TO_MEMORIZE_TABLE            = "flash_cards_to_memorize"
	FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE = "flash_cards_to_memorize_in_process"
)

func CountInDb(db sqlx.DB, boxId string, tableName string) int {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM `+tableName+` WHERE notion_data_base_id=?`, boxId)
	if err != nil {
		log.Fatalln(err)
	}
	return count
}

func CreateTableIfNotExist(db sqlx.DB, tableName string) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			page_id VARCHAR(255) NOT NULL PRIMARY KEY,
			image_url TEXT,
			notion_data_base_id TEXT,
			name TEXT,
			example TEXT,
			answer TEXT,
			know_level_1 TEXT,
			know_level_2 TEXT,
			know_level_3 TEXT,
			know_level_4 TEXT,
			know_level_5 TEXT,
			know_level_6 TEXT,
			know_level_7 TEXT,
			know_level_8 TEXT,
			know_level_9 TEXT,
			know_level_10 TEXT,
			know_level_11 TEXT,
			know_level_12 TEXT,
			know_level_13 TEXT,
			know_level_14 TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func NewFlashcardsFromDb(db sqlx.DB, tableName string) []Flashcard {
	query := `SELECT * FROM ` + tableName
	flashcards := []Flashcard{}
	err := db.Select(&flashcards, query)
	if err != nil {
		log.Fatal(err)
	}
	return flashcards
}

func (flashcard Flashcard) UpdateImage(db sqlx.DB, tableName string, newImage string) {
	query := `UPDATE ` + tableName + ` SET image_url=? WHERE page_id =?;`
	_, err := db.Exec(query, newImage, flashcard.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllFromBdByBoxId(db sqlx.DB, boxId string, tableName string) []Flashcard {
	flashcards := []Flashcard{}
	err := db.Select(&flashcards, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln(err)
	}
	return flashcards
}

func GetFormBox(db sqlx.DB, boxId string, tableName string) Flashcard {
	flashcard := Flashcard{}
	err := db.Get(&flashcard, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln(err)
	}
	return flashcard
}

func GetFromDdById(db sqlx.DB, id string, tableName string) Flashcard {
	flashcard := Flashcard{}
	err := db.Get(&flashcard, "SELECT * FROM "+tableName+" WHERE page_id=?", id)
	if err != nil {
		log.Fatalln(err)
	}
	return flashcard
}

func (flashcard Flashcard) RemoveFromDb(db sqlx.DB, tableName string) {
	_, err := db.Exec(`DELETE FROM `+tableName+` WHERE page_id = ?`, flashcard.Id)
	if err != nil {
		log.Fatalln(err)
	}
}

func InsertFlashCardsIntoDB(db sqlx.DB, flashCards []Flashcard, tableName string) {
	query := `INSERT INTO ` + tableName + ` (
		page_id, 
		image_url, 
		notion_data_base_id,
		name, 
		example, 
		answer,
		know_level_1,
		know_level_2,
		know_level_3,
		know_level_4, 
		know_level_5,
		know_level_6, 
		know_level_7, 
		know_level_8,
		know_level_9, 
		know_level_10, 
		know_level_11,
		know_level_12, 
		know_level_13,
		know_level_14
	) VALUES (
		:page_id, 
		:image_url, 
		:notion_data_base_id,
		:name,
		:example,
		:answer,
		:know_level_1,
		:know_level_2,
		:know_level_3,
		:know_level_4,
		:know_level_5,
		:know_level_6,
		:know_level_7,
		:know_level_8,
		:know_level_9,
		:know_level_10,
		:know_level_11,
		:know_level_12,
		:know_level_13,
		:know_level_14
	 );`
	_, err := db.NamedExec(query, flashCards)
	if err != nil {
		log.Fatal(err)
	}
}

func ClearFlashCardTable(db sqlx.DB, tableName string) {
	_, err := db.Exec(`DELETE FROM ` + tableName)
	if err != nil {
		log.Fatal(err)
	}
}
