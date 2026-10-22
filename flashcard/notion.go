package flashcard

import (
	"log"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/outbox"
	"spacedrepetitiongo/utils"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	notionApi "github.com/jomei/notionapi"
)

func NewFlashCards(pages []notionApi.Page) []Flashcard {
	flashCards := []Flashcard{}
	for _, page := range pages {
		flashCard := NewFlashCard(page)
		flashCards = append(flashCards, flashCard)
	}

	return flashCards
}

func NewFlashCard(page notionApi.Page) Flashcard {
	image := parseImage(page)
	return Flashcard{
		Id:          page.ID.String(),
		Image:       image,
		ImageSource: image,
		BoxId:       parseDbId(page.Parent.DatabaseID.String()),
		Name:        parseName(page),
		Example:     parseExample(page),
		Explanation: parseExplanation(page),
		KnowLevels:  parseKnowLevels(page),
	}
}

// NotionUpdate is the set of Notion properties that records this card's
// current know levels.
func (flashcard Flashcard) NotionUpdate() outbox.Update {
	update := outbox.Update{Checkboxes: map[string]bool{}}

	for level, value := range flashcard.KnowLevels {
		if value != nil {
			update.Checkboxes[KnowLevelProprtyName(level)] = *value
		}
	}

	if flashcard.BoxId == config.EnglishVocabularyId() {
		update.Numbers = map[string]float64{"Interval 14": 560}
	}

	return update
}

// EnqueueNotionUpdate durably queues this card's know levels for Notion.
// Call it before changing the local cache so the answer can't be lost.
func (flashcard Flashcard) EnqueueNotionUpdate(db sqlx.DB) Flashcard {
	if err := outbox.Enqueue(db, flashcard.Id, flashcard.NotionUpdate()); err != nil {
		log.Panic(err)
	}
	return flashcard
}

func KnowLevelProprtyName(level int) string {
	return "Know Level " + strconv.Itoa(level)
}

func parseImage(page notionApi.Page) *string {
	cover := page.Cover
	if cover == nil {
		return nil
	}
	externalCover := cover.External
	fileCover := cover.File
	if externalCover != nil {
		return &externalCover.URL
	}
	if fileCover != nil {
		return &fileCover.URL
	}
	return nil
}

func parseName(page notionApi.Page) string {
	nameProperty := page.Properties["Name"]
	titleProperty, _ := nameProperty.(*notionApi.TitleProperty)
	return utils.RichTextToString(titleProperty.Title)
}

func parseExample(page notionApi.Page) *string {
	exampleProperty := page.Properties["Example"]
	richTextProperty, ok := exampleProperty.(*notionApi.RichTextProperty)
	if !ok {
		return nil
	}
	str := utils.RichTextToString(richTextProperty.RichText)
	return &str
}

func parseExplanation(ppage notionApi.Page) *string {
	explanationProperty := ppage.Properties["Explanation"]
	explanationRichTextProperty, ok := explanationProperty.(*notionApi.RichTextProperty)
	var explanationStringBuffer strings.Builder
	if ok {
		explanationString := utils.RichTextToString(explanationRichTextProperty.RichText)
		if explanationString != "" {
			explanationStringBuffer.WriteString(explanationString)
		}
	}

	answerProperty := ppage.Properties["Answers"]
	answerRichTextProperty, ok := answerProperty.(*notionApi.RichTextProperty)
	if ok {
		answerString := utils.RichTextToString(answerRichTextProperty.RichText)
		if answerString != "" {
			if explanationStringBuffer.Len() > 0 {
				explanationStringBuffer.WriteString("\n")
			}
			explanationStringBuffer.WriteString(answerString)
		}
	}
	explanatuion := explanationStringBuffer.String()
	return &explanatuion
}

func parseKnowLevels(page notionApi.Page) map[int]*bool {
	knowLevels := make(map[int]*bool)
	for i := 1; i <= 14; i++ {
		knowLevels[i] = parseKnowLevel(i, page)
	}
	return knowLevels
}

func parseKnowLevel(level int, page notionApi.Page) *bool {
	knowLevelProperty := page.Properties[KnowLevelProprtyName(level)]
	checkboxProperty, ok := knowLevelProperty.(*notionApi.CheckboxProperty)
	if ok {
		return &checkboxProperty.Checkbox
	} else {
		return nil
	}
}

func parseDbId(id string) string {
	return strings.ReplaceAll(id, "-", "")
}
