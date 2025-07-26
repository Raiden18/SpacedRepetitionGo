package flashcard

import (
	"sort"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	"github.com/jmoiron/sqlx"
)

var (
	FORGOTTEN_VALUE = false
	MEMORIZED_VALUE = true
)

type Flashcard struct {
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

func NewMemorizingFlashcardFromDb(db sqlx.DB, id string) Flashcard {
	return GetFromDdById(db, id, FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
}

func NewRevisingFlashcardcFromDbById(db sqlx.DB, id string) Flashcard {
	return GetFromDdById(db, id, FLASH_CARDS_TO_REVISE_TABLE)
}

func NewRevisingFlashcardFromDbByBoxId(db sqlx.DB, boxId string) *Flashcard {
	if CountInDb(db, boxId, FLASH_CARDS_TO_REVISE_TABLE) > 0 {
		flashcard := GetFormBox(db, boxId, FLASH_CARDS_TO_REVISE_TABLE)
		return &flashcard
	} else {
		return nil
	}
}

func NewMemorizingFlashcardFromDbByBoxId(db sqlx.DB, boxId string) *Flashcard {
	if CountInDb(db, boxId, FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE) > 0 {
		flashcard := GetFormBox(db, boxId, FLASH_CARDS_TO_MEMORIZE_IN_PROCESS_TABLE)
		return &flashcard
	} else {
		return nil
	}
}

func (flashcard Flashcard) RemoveFrom(db sqlx.DB, tableName string) Flashcard {
	flashcard.RemoveFromDb(db, tableName)
	return flashcard
}

func (flashcard Flashcard) UpdateOnNotion(client notion.Client) Flashcard {
	go func() { flashcard.UpdatePageOnNotion(client) }()
	return flashcard
}

func (flashcard Flashcard) RemoveFromChat(bot telegram.Bot, id int) Flashcard {
	bot.DeleteMessage(id)
	return flashcard
}

func (flashcard Flashcard) InsertInto(db sqlx.DB, tableName string) Flashcard {
	InsertFlashCardsIntoDB(db, []Flashcard{flashcard}, tableName)
	return flashcard
}

func (flashcard *Flashcard) ToTelegramMessageToRevise() *FlashcardTelegramMessageToRevise {
	if flashcard == nil {
		return nil
	}
	message := NewFlashcardTelegramMessageToRevise(*flashcard)
	return &message
}

func (flashcard *Flashcard) ToTelegramMessageToMemorize() *FlashcardTelegramMessageToMemorize {
	if flashcard == nil {
		return nil
	}
	message := NewFlashcardTelegramMessageToMemorize(*flashcard)
	return &message
}

func (flashcard Flashcard) HasExplanation() bool {
	return flashcard.Explanation != nil && *flashcard.Explanation != ""
}

func (flashcard Flashcard) HasExample() bool {
	return flashcard.Example != nil && *flashcard.Example != ""
}

func (flashcard Flashcard) HasImage() bool {
	return flashcard.Image != nil && *flashcard.Image != ""
}

func (flashcard Flashcard) Recall() Flashcard {
	knowLevels := map[int]*bool{
		1:  flashcard.KnowLevel1,
		2:  flashcard.KnowLevel2,
		3:  flashcard.KnowLevel3,
		4:  flashcard.KnowLevel4,
		5:  flashcard.KnowLevel5,
		6:  flashcard.KnowLevel6,
		7:  flashcard.KnowLevel7,
		8:  flashcard.KnowLevel8,
		9:  flashcard.KnowLevel9,
		10: flashcard.KnowLevel10,
		11: flashcard.KnowLevel11,
		12: flashcard.KnowLevel12,
		13: flashcard.KnowLevel13,
		14: flashcard.KnowLevel14,
	}

	RecallAsMap(knowLevels)

	flashcard.KnowLevel1 = knowLevels[1]
	flashcard.KnowLevel2 = knowLevels[2]
	flashcard.KnowLevel3 = knowLevels[3]
	flashcard.KnowLevel4 = knowLevels[4]
	flashcard.KnowLevel5 = knowLevels[5]
	flashcard.KnowLevel6 = knowLevels[6]
	flashcard.KnowLevel7 = knowLevels[7]
	flashcard.KnowLevel8 = knowLevels[8]
	flashcard.KnowLevel9 = knowLevels[9]
	flashcard.KnowLevel10 = knowLevels[10]
	flashcard.KnowLevel11 = knowLevels[11]
	flashcard.KnowLevel12 = knowLevels[12]
	flashcard.KnowLevel13 = knowLevels[13]
	flashcard.KnowLevel14 = knowLevels[14]

	return flashcard
}

func (flashcard Flashcard) Forget() Flashcard {
	knowLevels := map[int]*bool{
		1:  flashcard.KnowLevel1,
		2:  flashcard.KnowLevel2,
		3:  flashcard.KnowLevel3,
		4:  flashcard.KnowLevel4,
		5:  flashcard.KnowLevel5,
		6:  flashcard.KnowLevel6,
		7:  flashcard.KnowLevel7,
		8:  flashcard.KnowLevel8,
		9:  flashcard.KnowLevel9,
		10: flashcard.KnowLevel10,
		11: flashcard.KnowLevel11,
		12: flashcard.KnowLevel12,
		13: flashcard.KnowLevel13,
		14: flashcard.KnowLevel14,
	}

	ForgetAsMap(knowLevels)

	flashcard.KnowLevel1 = knowLevels[1]
	flashcard.KnowLevel2 = knowLevels[2]
	flashcard.KnowLevel3 = knowLevels[3]
	flashcard.KnowLevel4 = knowLevels[4]
	flashcard.KnowLevel5 = knowLevels[5]
	flashcard.KnowLevel6 = knowLevels[6]
	flashcard.KnowLevel7 = knowLevels[7]
	flashcard.KnowLevel8 = knowLevels[8]
	flashcard.KnowLevel9 = knowLevels[9]
	flashcard.KnowLevel10 = knowLevels[10]
	flashcard.KnowLevel11 = knowLevels[11]
	flashcard.KnowLevel12 = knowLevels[12]
	flashcard.KnowLevel13 = knowLevels[13]
	flashcard.KnowLevel14 = knowLevels[14]

	return flashcard
}

func (flashcard Flashcard) Memorize() Flashcard {
	knowLevels := map[int]*bool{
		1:  flashcard.KnowLevel1,
		2:  flashcard.KnowLevel2,
		3:  flashcard.KnowLevel3,
		4:  flashcard.KnowLevel4,
		5:  flashcard.KnowLevel5,
		6:  flashcard.KnowLevel6,
		7:  flashcard.KnowLevel7,
		8:  flashcard.KnowLevel8,
		9:  flashcard.KnowLevel9,
		10: flashcard.KnowLevel10,
		11: flashcard.KnowLevel11,
		12: flashcard.KnowLevel12,
		13: flashcard.KnowLevel13,
		14: flashcard.KnowLevel14,
	}

	MemorizeAsMap(knowLevels)

	flashcard.KnowLevel1 = knowLevels[1]
	flashcard.KnowLevel2 = knowLevels[2]
	flashcard.KnowLevel3 = knowLevels[3]
	flashcard.KnowLevel4 = knowLevels[4]
	flashcard.KnowLevel5 = knowLevels[5]
	flashcard.KnowLevel6 = knowLevels[6]
	flashcard.KnowLevel7 = knowLevels[7]
	flashcard.KnowLevel8 = knowLevels[8]
	flashcard.KnowLevel9 = knowLevels[9]
	flashcard.KnowLevel10 = knowLevels[10]
	flashcard.KnowLevel11 = knowLevels[11]
	flashcard.KnowLevel12 = knowLevels[12]
	flashcard.KnowLevel13 = knowLevels[13]
	flashcard.KnowLevel14 = knowLevels[14]

	return flashcard
}

func ForgetAsMap(knowLevels map[int]*bool) {
	for key, value := range knowLevels {
		if value != nil {
			knowLevels[key] = &FORGOTTEN_VALUE
		}
	}
}

func MemorizeAsMap(knowLevels map[int]*bool) {
	ForgetAsMap(knowLevels)
	knowLevels[1] = &MEMORIZED_VALUE
}

func RecallAsMap(knowLevels map[int]*bool) {

	keys := make([]int, 0, len(knowLevels))
	for k := range knowLevels {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		value := knowLevels[key]
		if value != nil && !*value {
			knowLevels[key] = &MEMORIZED_VALUE
			break
		}
	}
}
