package notion

import (
	"context"
	"log"
	"spacedrepetitiongo/config"

	notionApi "github.com/jomei/notionapi"
)

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

func (client Client) FetchDataBase(id string) *notionApi.Database {
	db, error := client.ClientImpl.Database.Get(
		context.Background(),
		notionApi.DatabaseID(id),
	)
	if error != nil {
		log.Fatalln("Could not get Notiod Database: "+id, error)
	}

	return db
}

func (client Client) FetchPagesFromDb(id string, request *notionApi.DatabaseQueryRequest) []notionApi.Page {
	pages, err := client.ClientImpl.Database.Query(
		context.Background(),
		notionApi.DatabaseID(id),
		request,
	)
	if err != nil {
		log.Fatalln("Could not get pages from DB: "+id, err)
	}
	return pages.Results
}

func (client Client) UpdatePage(id string, updateRequest *notionApi.PageUpdateRequest) {
	pageId := notionApi.PageID(id)

	_, error := client.ClientImpl.Page.Update(context.Background(), pageId, updateRequest)
	if error != nil {
		log.Fatalln("Could not update Page", error)
	}
}
