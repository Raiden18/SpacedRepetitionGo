package notion

import (
	notionApi "github.com/jomei/notionapi"
)

func NewDatabaseQueryRequest(
	filters notionApi.Filter,
) notionApi.DatabaseQueryRequest {
	return notionApi.DatabaseQueryRequest{
		Filter: filters,
	}
}

func NewUpdateRequest(properties notionApi.Properties) notionApi.PageUpdateRequest {
	return notionApi.PageUpdateRequest{
		Properties: properties,
	}
}
