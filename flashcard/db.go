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
	Next        *string `db:"next"`
	Previous    *string `db:"previous"`
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
		Next:        flashcard.Next,
		Previous:    flashcard.Previous,
		KnowLevel1:  flashcard.KnowLevels[1],
		KnowLevel2:  flashcard.KnowLevels[2],
		KnowLevel3:  flashcard.KnowLevels[3],
		KnowLevel4:  flashcard.KnowLevels[4],
		KnowLevel5:  flashcard.KnowLevels[5],
		KnowLevel6:  flashcard.KnowLevels[6],
		KnowLevel7:  flashcard.KnowLevels[7],
		KnowLevel8:  flashcard.KnowLevels[8],
		KnowLevel9:  flashcard.KnowLevels[9],
		KnowLevel10: flashcard.KnowLevels[10],
		KnowLevel11: flashcard.KnowLevels[11],
		KnowLevel12: flashcard.KnowLevels[12],
		KnowLevel13: flashcard.KnowLevels[13],
		KnowLevel14: flashcard.KnowLevels[14],
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
		Next:        entity.Next,
		Previous:    entity.Previous,
		KnowLevels: map[int]*bool{
			1:  entity.KnowLevel1,
			2:  entity.KnowLevel2,
			3:  entity.KnowLevel3,
			4:  entity.KnowLevel4,
			5:  entity.KnowLevel5,
			6:  entity.KnowLevel6,
			7:  entity.KnowLevel7,
			8:  entity.KnowLevel8,
			9:  entity.KnowLevel9,
			10: entity.KnowLevel10,
			11: entity.KnowLevel11,
			12: entity.KnowLevel12,
			13: entity.KnowLevel13,
			14: entity.KnowLevel14,
		},
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
			next TEXT,
			previous TEXT,
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
		log.Fatal("Could not create New from DB. TableName="+tableName, err)
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

func (flashcard Flashcard) UpdateLinkedFlashCards(db sqlx.DB, tableName string) {
	query := `UPDATE ` + tableName + ` SET previous=?, next=? WHERE page_id =?;`
	_, err := db.Exec(query, flashcard.Previous, flashcard.Next, flashcard.Id)
	if err != nil {
		log.Fatal("Could not update linked flash cards", err)
	}
}

func GetAllFromBdByBoxId(db sqlx.DB, boxId string, tableName string) []Flashcard {
	entities := []FlashcardEntity{}
	err := db.Select(&entities, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln("Could not get all from BD by BoxId. TableName="+tableName, err)
	}
	return utils.Map(entities, toFlashcard)
}

func GetFormBox(db sqlx.DB, boxId string, tableName string) Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=?", boxId)
	if err != nil {
		log.Fatalln("Could not get from box. TableName="+tableName, err)
	}
	return toFlashcard(entity)
}

func GetLast(db sqlx.DB, boxId string, tableName string) *Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=? AND next IS NULL", boxId)
	if err != nil {
		return nil
	}
	flashcard := toFlashcard(entity)
	return &flashcard
}

func GetFirst(db sqlx.DB, boxId string, tableName string) *Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE notion_data_base_id=? AND previous IS NULL", boxId)
	if err != nil {
		return nil
	}
	flashCard := toFlashcard(entity)
	return &flashCard
}

func GetNextFromDb(db sqlx.DB, currentFlashcard Flashcard, tableName string) *Flashcard {
	if currentFlashcard.Next == nil {
		return nil
	}
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE next=?", *currentFlashcard.Next)
	if err != nil {
		log.Fatalln("Could not Get Next. TableName="+tableName, err)
	}
	nextFlashcard := toFlashcard(entity)
	return &nextFlashcard
}

func GetFromDdById(db sqlx.DB, id string, tableName string) *Flashcard {
	entity := FlashcardEntity{}
	err := db.Get(&entity, "SELECT * FROM "+tableName+" WHERE page_id=?", id)
	if err != nil {
		log.Println("Could not get by id. TableName="+tableName, err)
		return nil
	}
	flashCard := toFlashcard(entity)
	return &flashCard
}

func (flashcard Flashcard) RemoveFromDb(db sqlx.DB, tableName string) {
	_, err := db.Exec(`DELETE FROM `+tableName+` WHERE page_id = ?`, flashcard.Id)
	if err != nil {
		log.Fatalln("Could not deleted From Table. TableName="+tableName, err)
	}
}

func InsertIntoDB(db sqlx.DB, flashCards []Flashcard, tableName string) {
	entities := utils.Map(flashCards, toEntity)
	query := `INSERT IGNORE INTO ` + tableName + ` (
		page_id, 
		image_url, 
		notion_data_base_id,
		name, 
		example, 
		answer,
		next,
		previous,
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
		:next,
		:previous,
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
		log.Fatal("Could not insert into DB. TableName="+tableName, err)
	}
}

func ClearTable(db sqlx.DB, tableName string) {
	_, err := db.Exec(`DELETE FROM ` + tableName)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllIds(db sqlx.DB, tableName string) []string {
	ids := []string{}
	err := db.Select(&ids, "SELECT page_id FROM "+tableName)
	if err != nil {
		log.Fatalln("Could not get all ids. TableName="+tableName, err)
	}
	return ids
}

func DeleteMissing(db sqlx.DB, tableName string, existingIds map[string]struct{}) {
	ids := GetAllIds(db, tableName)
	for _, id := range ids {
		if _, ok := existingIds[id]; ok {
			continue
		}
		_, err := db.Exec(`DELETE FROM `+tableName+` WHERE page_id = ?`, id)
		if err != nil {
			log.Fatalln("Could not delete missing id. TableName="+tableName, err)
		}
	}
}
