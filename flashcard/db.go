package flashcard

import (
	"log"
	"spacedrepetitiongo/utils"

	"github.com/jmoiron/sqlx"
)

var (
	FLASH_CARDS_TO_REVISE_TABLE              = "flash_cards_to_revise"
	FLASH_CARDS_TO_MEMORIZE_TABLE            = "flash_cards_to_memorize"
	FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE = "flash_cards_to_memorize_in_process"
)

type FlashcardEntity struct {
	Id          string  `db:"page_id"`
	Image       *string `db:"image_url"`
	BoxId       string  `db:"notion_data_base_id"`
	Name        string  `db:"name"`
	Example     *string `db:"example"`
	Explanation *string `db:"answer"`
	KnowLevel1  *bool   `db:"know_level_1"`
	KnowLevel2  *bool   `db:"know_level_2"`
	KnowLevel3  *bool   `db:"know_level_3"`
	KnowLevel4  *bool   `db:"know_level_4"`
	KnowLevel5  *bool   `db:"know_level_5"`
	KnowLevel6  *bool   `db:"know_level_6"`
	KnowLevel7  *bool   `db:"know_level_7"`
	KnowLevel8  *bool   `db:"know_level_8"`
	KnowLevel9  *bool   `db:"know_level_9"`
	KnowLevel10 *bool   `db:"know_level_10"`
	KnowLevel11 *bool   `db:"know_level_11"`
	KnowLevel12 *bool   `db:"know_level_12"`
	KnowLevel13 *bool   `db:"know_level_13"`
	KnowLevel14 *bool   `db:"know_level_14"`
}

func toEntity(flashcard Flashcard) FlashcardEntity {
	return FlashcardEntity{
		Id:          flashcard.Id,
		Image:       flashcard.Image,
		BoxId:       flashcard.BoxId,
		Name:        flashcard.Name,
		Example:     flashcard.Example,
		Explanation: flashcard.Explanation,
		KnowLevel1:  flashcard.GetKnowLevels()[1],
		KnowLevel2:  flashcard.GetKnowLevels()[2],
		KnowLevel3:  flashcard.GetKnowLevels()[3],
		KnowLevel4:  flashcard.GetKnowLevels()[4],
		KnowLevel5:  flashcard.GetKnowLevels()[5],
		KnowLevel6:  flashcard.GetKnowLevels()[6],
		KnowLevel7:  flashcard.GetKnowLevels()[7],
		KnowLevel8:  flashcard.GetKnowLevels()[8],
		KnowLevel9:  flashcard.GetKnowLevels()[9],
		KnowLevel10: flashcard.GetKnowLevels()[10],
		KnowLevel11: flashcard.GetKnowLevels()[11],
		KnowLevel12: flashcard.GetKnowLevels()[12],
		KnowLevel13: flashcard.GetKnowLevels()[13],
		KnowLevel14: flashcard.GetKnowLevels()[14],
	}
}

func toFlashcard(entity FlashcardEntity) Flashcard {
	return Flashcard{
		Id:          entity.Id,
		Image:       entity.Image,
		BoxId:       entity.BoxId,
		Name:        entity.Name,
		Example:     entity.Example,
		Explanation: entity.Explanation,
		KnowLevel1:  entity.KnowLevel1,
		KnowLevel2:  entity.KnowLevel2,
		KnowLevel3:  entity.KnowLevel3,
		KnowLevel4:  entity.KnowLevel4,
		KnowLevel5:  entity.KnowLevel5,
		KnowLevel6:  entity.KnowLevel6,
		KnowLevel7:  entity.KnowLevel7,
		KnowLevel8:  entity.KnowLevel8,
		KnowLevel9:  entity.KnowLevel9,
		KnowLevel10: entity.KnowLevel10,
		KnowLevel11: entity.KnowLevel11,
		KnowLevel12: entity.KnowLevel12,
		KnowLevel13: entity.KnowLevel13,
		KnowLevel14: entity.KnowLevel14,
	}
}

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
	entities := []FlashcardEntity{}
	err := db.Select(&entities, query)
	if err != nil {
		log.Fatal(err)
	}
	return utils.Map(entities, toFlashcard)
}

func (flashcard Flashcard) UpdateImage(db sqlx.DB, tableName string, newImage string) {
	query := `UPDATE ` + tableName + ` SET image_url=? WHERE page_id =?;`
	_, err := db.Exec(query, newImage, flashcard.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllFromBdByBoxId(db sqlx.DB, boxId string, tableName string) []Flashcard {
	entities := []FlashcardEntity{}
	err := db.Select(&entities, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln(err)
	}
	return utils.Map(entities, toFlashcard)
}

func GetFormBox(db sqlx.DB, boxId string, tableName string) Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln(err)
	}
	return toFlashcard(entity)
}

func GetFromDdById(db sqlx.DB, id string, tableName string) Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE page_id=?", id)
	if err != nil {
		log.Fatalln(err)
	}
	return toFlashcard(entity)
}

func (flashcard Flashcard) RemoveFromDb(db sqlx.DB, tableName string) {
	_, err := db.Exec(`DELETE FROM `+tableName+` WHERE page_id = ?`, flashcard.Id)
	if err != nil {
		log.Fatalln(err)
	}
}

func InsertFlashCardsIntoDB(db sqlx.DB, flashCards []Flashcard, tableName string) {
	entities := utils.Map(flashCards, toEntity)
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
	_, err := db.NamedExec(query, entities)
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
