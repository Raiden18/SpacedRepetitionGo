package box

import (
	"spacedrepetitiongo/utils"
	"strings"

	notionApi "github.com/jomei/notionapi"
)

func NewBoxes(databases []notionApi.Database) []Box {
	boxes := []Box{}
	for _, database := range databases {
		box := NewBox(database)
		boxes = append(boxes, box)
	}
	return boxes
}

func NewBox(database notionApi.Database) Box {
	return Box{
		Id:   parseNotionDbId(database),
		Name: utils.RichTextToString(database.Title),
	}
}

func parseNotionDbId(database notionApi.Database) string {
	return RemoveDashFromId(
		database.ID.String(),
	)
}

func RemoveDashFromId(id string) string {
	return strings.ReplaceAll(id, "-", "")
}
