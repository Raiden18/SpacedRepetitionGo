package flashcard

import (
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/utils"
	"strconv"
	"strings"

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
	return Flashcard{
		Id:          page.ID.String(),
		Image:       parseImage(page),
		BoxId:       parseDbId(page.Parent.DatabaseID.String()),
		Name:        parseName(page),
		Example:     parseExample(page),
		Explanation: parseExplanation(page),
		KnowLevels:  parseKnowLevels(page),
	}
}

func (flashcard Flashcard) UpdatePageOnNotion(client notion.Client) {
	properties := notionApi.Properties{}

	for level, value := range flashcard.KnowLevels {
		if value != nil {
			properties[KnowLevelProprtyName(level)] = notionApi.CheckboxProperty{Checkbox: *value}
		}
	}

	if flashcard.BoxId == config.EnglishVocabularyId() {
		properties["Interval 14"] = notionApi.NumberProperty{Number: 560}
	}

	updateRequest := notion.NewUpdateRequest(properties)
	client.UpdatePage(flashcard.Id, &updateRequest)
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
