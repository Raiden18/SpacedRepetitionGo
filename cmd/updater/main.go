package main

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/utils"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	notionApi "github.com/jomei/notionapi"
)

func main() {
	db := utils.OpenDb()
	notionClient := notion.NewClient()
	observedDatabasesIds := config.GetObservedDatabasesIds()
	boxes := fetchFromNotion(observedDatabasesIds, notionClient)
	reviseFlashcardsRequest := notion.NewDatabaseQueryRequest(
		notion.And(
			notion.KnowLevel(1, true),
			notion.Show(true),
		),
	)
	memorizeFlashcardsRequest := notion.NewDatabaseQueryRequest(
		notion.And(
			notion.KnowLevel(1, false),
			notion.KnowLevel(2, false),
			notion.KnowLevel(3, false),
			notion.KnowLevel(4, false),
			notion.Show(true),
		),
	)

	pagesToRevise := fetchFlashCards(boxes, notionClient, &reviseFlashcardsRequest)
	pagesToMemorize := fetchFlashCards(boxes, notionClient, &memorizeFlashcardsRequest)

	flashCardsToRevise := flashcard.NewFlashCards(pagesToRevise)
	flashCardsToMemorize := flashcard.NewFlashCards(pagesToMemorize)

	box.ClearTable(db)
	box.InsertIntoDB(db, boxes)

	flashcard.ClearFlashCardTable(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	insertFlashCards(db, flashCardsToRevise, flashcard.FLASH_CARDS_TO_REVISE_TABLE)

	flashcard.ClearFlashCardTable(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)
	insertFlashCards(db, flashCardsToMemorize, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)

	defer db.Close()
}

func fetchFlashCards(boxes []box.Box, client notion.Client, request *notionApi.DatabaseQueryRequest) []notionApi.Page {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		pages []notionApi.Page
	)

	for _, box_ := range boxes {
		wg.Add(1)
		go func(b box.Box) {
			defer wg.Done()
			dbPages := client.FetchPagesFromDb(b.Id, request)
			mutex.Lock()
			pages = append(pages, dbPages...)
			mutex.Unlock()
		}(box_)
	}

	wg.Wait()
	return pages
}

func fetchFromNotion(
	ids []string,
	client notion.Client,
) []box.Box {
	var (
		wg        sync.WaitGroup
		mutex     sync.Mutex
		databases []notionApi.Database
	)
	for _, databaseId := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			db := client.FetchDataBase(databaseId)
			mutex.Lock()
			databases = append(databases, *db)
			mutex.Unlock()
		}(databaseId)
	}
	wg.Wait()
	return box.NewBoxes(databases)
}

func insertFlashCards(db sqlx.DB, flashCards []flashcard.Flashcard, tableName string) {
	if len(flashCards) > 0 {
		flashcard.InsertFlashCardsIntoDB(db, flashCards, tableName)
	}
}
