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
		KnowLevel1:  parseKnowLevel(1, page),
		KnowLevel2:  parseKnowLevel(2, page),
		KnowLevel3:  parseKnowLevel(3, page),
		KnowLevel4:  parseKnowLevel(4, page),
		KnowLevel5:  parseKnowLevel(5, page),
		KnowLevel6:  parseKnowLevel(6, page),
		KnowLevel7:  parseKnowLevel(7, page),
		KnowLevel8:  parseKnowLevel(8, page),
		KnowLevel9:  parseKnowLevel(9, page),
		KnowLevel10: parseKnowLevel(10, page),
		KnowLevel11: parseKnowLevel(11, page),
		KnowLevel12: parseKnowLevel(12, page),
		KnowLevel13: parseKnowLevel(13, page),
		KnowLevel14: parseKnowLevel(14, page),
	}
}

func (flashcard Flashcard) UpdatePageOnNotion(client notion.Client) {
	knowLevelProperties := notionApi.Properties{}

	newKnowLevelProperty(flashcard.KnowLevel1, knowLevelProperties, 1, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel2, knowLevelProperties, 2, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel3, knowLevelProperties, 3, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel4, knowLevelProperties, 4, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel5, knowLevelProperties, 5, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel6, knowLevelProperties, 6, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel7, knowLevelProperties, 7, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel8, knowLevelProperties, 8, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel9, knowLevelProperties, 9, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel10, knowLevelProperties, 10, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel11, knowLevelProperties, 11, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel12, knowLevelProperties, 12, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel13, knowLevelProperties, 13, flashcard)
	newKnowLevelProperty(flashcard.KnowLevel14, knowLevelProperties, 14, flashcard)

	updateRequest := notion.NewUpdateRequest(knowLevelProperties)
	client.UpdatePage(flashcard.Id, &updateRequest)
}

func newKnowLevelProperty(knowLevel *bool, properties notionApi.Properties, level int, flashcard Flashcard) {
	if knowLevel != nil {
		propertyName := "Know Level " + strconv.Itoa(level)
		properties[propertyName] = notionApi.CheckboxProperty{Checkbox: *knowLevel}
	}
	if flashcard.BoxId == config.EnglishVocabularyId() {
		properties["Interval 14"] = notionApi.NumberProperty{Number: 560}
	}
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
	return titleProperty.Title[0].Text.Content
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
	var explanationString string = ""
	if ok {
		explanationString = utils.RichTextToString(explanationRichTextProperty.RichText)
	}

	answerProperty := ppage.Properties["Answers"]
	answerRichTextProperty, ok := answerProperty.(*notionApi.RichTextProperty)
	var answerString string = ""
	if ok {
		answerString = utils.RichTextToString(answerRichTextProperty.RichText)
	}
	str := answerString + "\n" + explanationString
	return &str
}

func parseKnowLevel(level int, page notionApi.Page) *bool {
	knowLevelProperty := page.Properties["Know Level "+strconv.Itoa(level)]
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
