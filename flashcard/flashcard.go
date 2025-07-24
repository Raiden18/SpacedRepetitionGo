package flashcard

import (
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	"github.com/jmoiron/sqlx"
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
	recalledValue := true

	levels := []*bool{
		flashcard.KnowLevel1,
		flashcard.KnowLevel2,
		flashcard.KnowLevel3,
		flashcard.KnowLevel4,
		flashcard.KnowLevel5,
		flashcard.KnowLevel6,
		flashcard.KnowLevel7,
		flashcard.KnowLevel8,
		flashcard.KnowLevel9,
		flashcard.KnowLevel10,
		flashcard.KnowLevel11,
		flashcard.KnowLevel12,
		flashcard.KnowLevel13,
		flashcard.KnowLevel14,
	}

	fieldPtrs := []**bool{
		&flashcard.KnowLevel1,
		&flashcard.KnowLevel2,
		&flashcard.KnowLevel3,
		&flashcard.KnowLevel4,
		&flashcard.KnowLevel5,
		&flashcard.KnowLevel6,
		&flashcard.KnowLevel7,
		&flashcard.KnowLevel8,
		&flashcard.KnowLevel9,
		&flashcard.KnowLevel10,
		&flashcard.KnowLevel11,
		&flashcard.KnowLevel12,
		&flashcard.KnowLevel13,
		&flashcard.KnowLevel14,
	}

	for i := 0; i < len(levels)-1; i++ {
		curr := levels[i]
		next := levels[i+1]

		if curr != nil && *curr && next != nil && !*next {
			*fieldPtrs[i+1] = &recalledValue
			break
		}
	}

	return flashcard
}

func (flashcard Flashcard) Forget() Flashcard {
	forgottenValue := false
	if flashcard.KnowLevel1 != nil {
		flashcard.KnowLevel1 = &forgottenValue
	}
	if flashcard.KnowLevel2 != nil {
		flashcard.KnowLevel2 = &forgottenValue
	}
	if flashcard.KnowLevel3 != nil {
		flashcard.KnowLevel3 = &forgottenValue
	}
	if flashcard.KnowLevel4 != nil {
		flashcard.KnowLevel4 = &forgottenValue
	}
	if flashcard.KnowLevel5 != nil {
		flashcard.KnowLevel5 = &forgottenValue
	}
	if flashcard.KnowLevel6 != nil {
		flashcard.KnowLevel6 = &forgottenValue
	}
	if flashcard.KnowLevel7 != nil {
		flashcard.KnowLevel7 = &forgottenValue
	}
	if flashcard.KnowLevel8 != nil {
		flashcard.KnowLevel8 = &forgottenValue
	}
	if flashcard.KnowLevel9 != nil {
		flashcard.KnowLevel9 = &forgottenValue
	}
	if flashcard.KnowLevel10 != nil {
		flashcard.KnowLevel10 = &forgottenValue
	}
	if flashcard.KnowLevel11 != nil {
		flashcard.KnowLevel11 = &forgottenValue
	}
	if flashcard.KnowLevel12 != nil {
		flashcard.KnowLevel12 = &forgottenValue
	}
	if flashcard.KnowLevel13 != nil {
		flashcard.KnowLevel13 = &forgottenValue
	}
	if flashcard.KnowLevel14 != nil {
		flashcard.KnowLevel14 = &forgottenValue
	}
	return flashcard
}

func (flashcard Flashcard) Memorize() Flashcard {
	true_ := true
	false_ := false

	if flashcard.KnowLevel1 != nil {
		flashcard.KnowLevel1 = &true_
	}
	if flashcard.KnowLevel2 != nil {
		flashcard.KnowLevel2 = &false_
	}
	if flashcard.KnowLevel3 != nil {
		flashcard.KnowLevel3 = &false_
	}
	if flashcard.KnowLevel4 != nil {
		flashcard.KnowLevel4 = &false_
	}
	if flashcard.KnowLevel5 != nil {
		flashcard.KnowLevel5 = &false_
	}
	if flashcard.KnowLevel6 != nil {
		flashcard.KnowLevel6 = &false_
	}
	if flashcard.KnowLevel7 != nil {
		flashcard.KnowLevel7 = &false_
	}
	if flashcard.KnowLevel8 != nil {
		flashcard.KnowLevel8 = &false_
	}
	if flashcard.KnowLevel9 != nil {
		flashcard.KnowLevel9 = &false_
	}
	if flashcard.KnowLevel10 != nil {
		flashcard.KnowLevel10 = &false_
	}
	if flashcard.KnowLevel11 != nil {
		flashcard.KnowLevel11 = &false_
	}
	if flashcard.KnowLevel12 != nil {
		flashcard.KnowLevel12 = &false_
	}
	if flashcard.KnowLevel13 != nil {
		flashcard.KnowLevel13 = &false_
	}
	if flashcard.KnowLevel14 != nil {
		flashcard.KnowLevel14 = &false_
	}
	return flashcard
}
