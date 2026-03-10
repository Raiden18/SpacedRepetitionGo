package box

import (
	"spacedrepetitiongo/utils"

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
	nameProperty := page.Properties["Name"]
	titleProperty, ok := nameProperty.(*notionApi.TitleProperty)
	if !ok {
		return ""
	}
	return utils.RichTextToString(titleProperty.Title)
}

func parseObservedDatabaseId(page notionApi.Page) string {
	idProperty := page.Properties["Id"]
	richTextProperty, ok := idProperty.(*notionApi.RichTextProperty)
	if !ok {
		return ""
	}
	return utils.RichTextToString(richTextProperty.RichText)
}

func parseObservedDatabaseObservable(page notionApi.Page) bool {
	observableProperty := page.Properties["Observable"]
	checkboxProperty, ok := observableProperty.(*notionApi.CheckboxProperty)
	if !ok {
		return false
	}
	return checkboxProperty.Checkbox
}
