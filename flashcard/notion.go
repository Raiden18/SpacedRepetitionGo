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
	properties := notionApi.Properties{}

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

	for level, value := range knowLevels {
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
