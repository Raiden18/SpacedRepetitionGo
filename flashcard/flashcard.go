package flashcard

import (
	"context"
	"fmt"
	"sort"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/telegram"

	"github.com/jmoiron/sqlx"
	openai "github.com/sashabaranov/go-openai"
)

var (
	FORGOTTEN_VALUE = false
	MEMORIZED_VALUE = true
)

type Flashcard struct {
	Id          string
	Image       *string
	BoxId       string
	Name        string
	Example     *string
	Explanation *string
	KnowLevels  map[int]*bool
}

func GenerateFromGPT(
	language string,
	word string,
	openAi *openai.Client,
) Flashcard {
	resp, err := openAi.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a Greek language tutor. Always reply with exactly one complete sentence in Greek.",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(
						"Write exactly one complete sentence in Greek using the word \"%s\". "+
							"The sentence must be simple, natural, and suitable for someone learning Greek. "+
							"Do not add translations, explanations, or any text before or after the sentence.",
						word,
					),
				},
			},
			MaxTokens: 50,
		},
	)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
	}

	mesasges := ""

	for _, choice := range resp.Choices {
		mesasges += choice.Message.Content + "\n"
	}
	return Flashcard{
		Id:          "GPT_GENERATED",
		Image:       nil,
		BoxId:       "NO",
		Name:        word,
		Example:     &resp.Choices[0].Message.Content,
		Explanation: nil,
		KnowLevels:  make(map[int]*bool),
	}
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

func HasImage(flashcard Flashcard) bool {
	return flashcard.HasImage()
}

func (flashcard Flashcard) Recall() Flashcard {
	RecallAsMap(flashcard.KnowLevels)
	return flashcard
}

func (flashcard Flashcard) Forget() Flashcard {
	ForgetAsMap(flashcard.KnowLevels)
	return flashcard
}

func (flashcard Flashcard) Memorize() Flashcard {
	MemorizeAsMap(flashcard.KnowLevels)
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
