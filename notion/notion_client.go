package notion

import (
	"context"
	"fmt"
	"spacedrepetitiongo/config"

	notionApi "github.com/jomei/notionapi"
)

// Notion returns at most this many results per query page.
const maxPageSize = 100

type Client struct {
	ClientImpl *notionApi.Client
}

func NewClient() Client {
	return Client{
		ClientImpl: notionApi.NewClient(
			notionApi.Token(
				config.NotionApiKey(),
			),
		),
	}
}

func (client Client) FetchDataBase(id string) (*notionApi.Database, error) {
	db, err := client.ClientImpl.Database.Get(
		context.Background(),
		notionApi.DatabaseID(id),
	)
	if err != nil {
		return nil, fmt.Errorf("get Notion database %s: %w", id, err)
	}
	return db, nil
}

// FetchPagesFromDb returns every page matching the request, following
// Notion's pagination until has_more is false. It returns an error rather
// than a partial result, so callers never mistake a truncated list for the
// full one.
func (client Client) FetchPagesFromDb(id string, request *notionApi.DatabaseQueryRequest) ([]notionApi.Page, error) {
	// Copy so the caller's request (shared across goroutines) is not mutated.
	query := *request
	query.PageSize = maxPageSize
	query.StartCursor = ""

	pages := []notionApi.Page{}
	for {
		response, err := client.ClientImpl.Database.Query(
			context.Background(),
			notionApi.DatabaseID(id),
			&query,
		)
		if err != nil {
			return nil, fmt.Errorf("query Notion database %s: %w", id, err)
		}
		pages = append(pages, response.Results...)
		if !response.HasMore || response.NextCursor == "" {
			return pages, nil
		}
		query.StartCursor = response.NextCursor
	}
}

func (client Client) UpdatePage(id string, updateRequest *notionApi.PageUpdateRequest) error {
	_, err := client.ClientImpl.Page.Update(context.Background(), notionApi.PageID(id), updateRequest)
	if err != nil {
		return fmt.Errorf("update Notion page %s: %w", id, err)
	}
	return nil
}
