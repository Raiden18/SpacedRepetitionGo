package jobs

import (
	"log"
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/outbox"
	"spacedrepetitiongo/schema"
	"spacedrepetitiongo/utils"
	"sync"
	"time"

	notionApi "github.com/jomei/notionapi"
)

// Update syncs the local queues from Notion. Everything is fetched before
// anything is written, and any Notion error aborts the run, so a failed or
// partial fetch can never delete cards from the local cache.
func Update() {
	db := utils.OpenDb()
	defer db.Close()
	schema.Ensure(db)

	notionClient := notion.NewClient()

	// Send queued answers first, so Notion is as current as possible.
	outbox.Drain(db, notionClient, 30*time.Second)
	pending, err := outbox.PendingPageIds(db)
	if err != nil {
		log.Fatal(err)
	}

	observedDatabases, err := fetchObservedDatabases(config.GetObservedDatabasesId(), notionClient)
	if err != nil {
		log.Fatal("Aborting update, nothing changed: ", err)
	}
	boxes, err := fetchBoxes(observedDatabases, notionClient)
	if err != nil {
		log.Fatal("Aborting update, nothing changed: ", err)
	}

	reviseFlashcardsRequest := notion.NewDatabaseQueryRequest(
		notion.And(
			KnowLevel(1, true),
			Show(true),
		),
	)
	memorizeFlashcardsRequest := notion.NewDatabaseQueryRequest(
		notion.And(
			KnowLevel(1, false),
			KnowLevel(2, false),
			KnowLevel(3, false),
			KnowLevel(4, false),
			Show(true),
		),
	)

	flashCardsToRevise, err := fetchFlashCards(boxes, notionClient, &reviseFlashcardsRequest, pending)
	if err != nil {
		log.Fatal("Aborting update, nothing changed: ", err)
	}
	flashCardsToMemorize, err := fetchFlashCards(boxes, notionClient, &memorizeFlashcardsRequest, pending)
	if err != nil {
		log.Fatal("Aborting update, nothing changed: ", err)
	}

	box.InsertIntoDB(db, boxes)

	// Pages with answers still in the outbox are left exactly as they are
	// locally: Notion's copy of them is stale.
	storedImages := flashcard.StoredImages(db)
	flashcard.Sync(
		db,
		flashcard.KeepDownloadedImages(flashCardsToRevise, storedImages),
		flashcard.KeepDownloadedImages(flashCardsToMemorize, storedImages),
		pending,
	)
}

func fetchObservedDatabases(
	observedDatabasesId string,
	client notion.Client,
) ([]string, error) {
	request := notion.NewEmptyDatabaseQueryRequest()
	pages, err := client.FetchPagesFromDb(observedDatabasesId, &request)
	if err != nil {
		return nil, err
	}
	databases := box.NewObservedDatabases(pages)
	ids := []string{}
	for _, database := range databases {
		if database.IsObservable() && database.Id != "" {
			ids = append(ids, database.Id)
		}
	}

	return ids, nil
}

// fetchFlashCards fetches every matching card of every box, in Notion's order
// within each box. Cards with answers still waiting in the outbox are left
// out: Notion's copy of them is stale, and the local cache already reflects
// the answer.
func fetchFlashCards(
	boxes []box.Box,
	client notion.Client,
	request *notionApi.DatabaseQueryRequest,
	pending map[string]struct{},
) ([]flashcard.Flashcard, error) {
	var (
		wg               sync.WaitGroup
		mutex            sync.Mutex
		flashCardsResult []flashcard.Flashcard
		firstErr         error
	)

	for _, b := range boxes {
		wg.Add(1)
		go func(b box.Box) {
			defer wg.Done()
			dbPages, err := client.FetchPagesFromDb(b.Id, request)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			flashCards := utils.Filter(flashcard.NewFlashCards(dbPages), func(card flashcard.Flashcard) bool {
				_, isPending := pending[card.Id]
				return !isPending
			})
			flashCardsResult = append(flashCardsResult, flashCards...)
		}(b)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return flashCardsResult, nil
}

func fetchBoxes(
	ids []string,
	client notion.Client,
) ([]box.Box, error) {
	var (
		wg        sync.WaitGroup
		mutex     sync.Mutex
		databases []notionApi.Database
		firstErr  error
	)
	for _, databaseId := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			db, err := client.FetchDataBase(id)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			databases = append(databases, *db)
		}(databaseId)
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return box.NewBoxes(databases), nil
}

func KnowLevel(level int, equals bool) notionApi.PropertyFilter {
	return notion.PropertyCheckBox(flashcard.KnowLevelProprtyName(level), equals)
}

func Show(equals bool) notionApi.PropertyFilter {
	return notion.PropertyCheckBox("Show", equals)
}
