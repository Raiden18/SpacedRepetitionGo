package jobs

import (
	"spacedrepetitiongo/box"
	"spacedrepetitiongo/config"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/notion"
	"spacedrepetitiongo/utils"
	"sync"

	"github.com/jmoiron/sqlx"
	notionApi "github.com/jomei/notionapi"
)

func Update() {
	db := utils.OpenDb()
	notionClient := notion.NewClient()
	observedDatabasesIds := config.GetObservedDatabasesIds()
	boxes := fetchFromNotion(observedDatabasesIds, notionClient)
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

	flashCardsToRevise := fetchFlashCards(boxes, notionClient, &reviseFlashcardsRequest)
	flashCardsToMemorize := fetchFlashCards(boxes, notionClient, &memorizeFlashcardsRequest)

	box.InsertIntoDB(db, boxes)

	syncFlashCards(db, flashCardsToRevise, flashcard.FLASH_CARDS_TO_REVISE_TABLE)
	syncFlashCards(db, flashCardsToMemorize, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE)

	defer db.Close()
}

func fetchFlashCards(boxes []box.Box, client notion.Client, request *notionApi.DatabaseQueryRequest) []flashcard.Flashcard {
	var (
		wg               sync.WaitGroup
		mutex            sync.Mutex
		flashCardsResult []flashcard.Flashcard
	)

	for _, box_ := range boxes {
		wg.Add(1)
		go func(b box.Box) {
			defer wg.Done()
			dbPages := client.FetchPagesFromDb(b.Id, request)
			flashCards := flashcard.NewFlashCards(dbPages)
			orderedFlashCards := orderFlashCards(flashCards)
			mutex.Lock()
			flashCardsResult = append(flashCardsResult, orderedFlashCards...)
			mutex.Unlock()
		}(box_)
	}

	wg.Wait()
	return flashCardsResult
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

func orderFlashCards(flashCards []flashcard.Flashcard) []flashcard.Flashcard {
	orderedFlashCards := make([]flashcard.Flashcard, len(flashCards))
	for index, flashCard := range flashCards {
		if index == 0 {
			flashCard.Previous = nil
		} else {
			prevFlashCard := flashCards[index-1]
			flashCard.Previous = &prevFlashCard.Id
		}
		if index+1 < len(flashCards) {
			nextFlashCard := flashCards[index+1]
			flashCard.Next = &nextFlashCard.Id
		} else {
			flashCard.Next = nil
		}
		orderedFlashCards[index] = flashCard
	}
	return orderedFlashCards
}

func syncFlashCards(db sqlx.DB, flashCards []flashcard.Flashcard, tableName string) {
	notionIds := make(map[string]struct{}, len(flashCards))
	for _, card := range flashCards {
		notionIds[card.Id] = struct{}{}
	}

	flashcard.DeleteMissing(db, tableName, notionIds)
	insertFlashCards(db, flashCards, tableName)
}

func insertFlashCards(db sqlx.DB, flashCards []flashcard.Flashcard, tableName string) {
	if len(flashCards) == 0 {
		return
	}
	flashcard.InsertIntoDB(db, flashCards, tableName)
}

func KnowLevel(level int, equals bool) notionApi.PropertyFilter {
	return notion.PropertyCheckBox(flashcard.KnowLevelProprtyName(level), equals)
}

func Show(equals bool) notionApi.PropertyFilter {
	return notion.PropertyCheckBox("Show", equals)
}
