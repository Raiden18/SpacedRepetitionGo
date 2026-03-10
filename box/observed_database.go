package box

import (
	"spacedrepetitiongo/utils"
	"strings"

	notionApi "github.com/jomei/notionapi"
)

type ObservedDatabase struct {
	Name       string
	Id         string
	Observable bool
}

func NewObservedDatabases(pages []notionApi.Page) []ObservedDatabase {
	observedDatabases := []ObservedDatabase{}
	for _, page := range pages {
		observedDatabase := NewObservedDatabase(page)
		observedDatabases = append(observedDatabases, observedDatabase)
	}

	return observedDatabases
}

func NewObservedDatabase(page notionApi.Page) ObservedDatabase {
	return ObservedDatabase{
		Name:       parseObservedDatabaseName(page),
		Id:         parseObservedDatabaseId(page),
		Observable: parseObservedDatabaseObservable(page),
	}
}

func (database ObservedDatabase) IsObservable() bool {
	return database.Observable
}

func parseObservedDatabaseName(page notionApi.Page) string {
	for propertyName, propertyValue := range page.Properties {
		if !strings.EqualFold(propertyName, "name") {
			continue
		}
		titleProperty, ok := propertyValue.(*notionApi.TitleProperty)
		if !ok {
			return ""
		}
		return utils.RichTextToString(titleProperty.Title)
	}
	return ""
}

func parseObservedDatabaseId(page notionApi.Page) string {
	for propertyName, propertyValue := range page.Properties {
		if !strings.EqualFold(propertyName, "id") {
			continue
		}
		richTextProperty, ok := propertyValue.(*notionApi.RichTextProperty)
		if !ok {
			return ""
		}
		return utils.RichTextToString(richTextProperty.RichText)
	}
	return ""
}

func parseObservedDatabaseObservable(page notionApi.Page) bool {
	for propertyName, propertyValue := range page.Properties {
		if !strings.EqualFold(propertyName, "observable") {
			continue
		}
		checkboxProperty, ok := propertyValue.(*notionApi.CheckboxProperty)
		if !ok {
			return false
		}
		return checkboxProperty.Checkbox
	}
	return false
}
