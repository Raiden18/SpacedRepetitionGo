package flashcard

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/openai"
	"spacedrepetitiongo/telegram"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
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
	openAi *openai.OpenAiClient,
) Flashcard {
	englishMeaning := openAi.Ask(
		"You are a Greek language tutor",
		fmt.Sprintf(
			"Give english translation for the word \"%s\""+
				"The english must be simple, natural, and suitable"+
				"Do not add any text before or after the translation."+
				"Verbs must be in nominative for \"I\", nouns in singular form",
			word,
		),
	)

	urlImage := openAi.CreateImage(
		fmt.Sprintf(
			"Create an image of the word \"%s\" in clip art style, that represents the word visually. Must be has jpeg format.",
			word,
		),
	)

	log.Printf(urlImage)
	downloadImage(
		resty.New(),
		urlImage,
		"/root/repetition/images/",
		"gpt_generated",
	)

	url, _ := findFileByNameWithoutExt(
		"/root/repetition/images/",
		"gpt_generated",
	)
	return Flashcard{
		Id:          "GPT_GENERATED",
		Image:       &url,
		BoxId:       "NO",
		Name:        englishMeaning,
		Example:     nil,
		Explanation: &word,
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

func downloadImage(client *resty.Client, url, folder, baseFilename string) error {
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+config.OpenAiApiKey()).
		Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.RawBody().Close()

	ext := ".jpg" // Assuming OpenAI always returns JPEG

	fullPath := filepath.Join(folder, baseFilename+ext)

	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.RawBody())
	if err != nil {
		return fmt.Errorf("failed to copy image data to file: %w", err)
	}
	return nil
}

func findFileByNameWithoutExt(folderPath, baseName string) (string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		nameWithoutExt := strings.TrimSuffix(name, ext)

		if nameWithoutExt == baseName {
			return filepath.Join(folderPath, name), nil
		}
	}

	return "", fmt.Errorf("file %q not found in %q", baseName, folderPath)
}
