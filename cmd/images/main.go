package main

import (
	"log"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/utils"

	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db := utils.OpenDb()
	client := resty.New()
	imagesFolder := "/root/repetition/images/"
	createFolderIfNotExist(imagesFolder)
	deleteAllFilesFromFolder(imagesFolder)

	flashCardsToRevise := flashCardsWithImages(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE),
	)
	flashCardsToMemorize := flashCardsWithImages(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
	)

	allFlashCards := append(flashCardsToMemorize, flashCardsToRevise...)

	downloadImages(client, imagesFolder, allFlashCards)

	converts := map[string]func(path string){
		".jfif": convertJfifToJpg,
		".htm":  convertHtmtoPng,
		".svg":  convertSVGtoPNG,
	}

	findFilesAndConvert(imagesFolder, converts)

	updateImagesInDb(db, flashCardsToRevise, flashcard.FLASH_CARDS_TO_REVISE_TABLE, imagesFolder)
	updateImagesInDb(db, flashCardsToMemorize, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE, imagesFolder)
}

func flashCardsWithImages(flashcards []flashcard.Flashcard) []flashcard.Flashcard {
	flashcardsWithImages := []flashcard.Flashcard{}

	for _, flashcard := range flashcards {
		if flashcard.Image != nil {
			flashcardsWithImages = append(flashcardsWithImages, flashcard)
		}
	}

	return flashcardsWithImages
}

func downloadImages(client *resty.Client, folder string, flashcard []flashcard.Flashcard) {
	for _, flashcardWithOldImage := range flashcard {
		originalUrl := *flashcardWithOldImage.Image
		if strings.Contains(originalUrl, "base64") {
			convertBase64ToImage(originalUrl, folder, flashcardWithOldImage.Id)
		} else if strings.HasPrefix(originalUrl, "https://www.notion.so/image/") {
			decoded, _ := url.QueryUnescape(
				substringBefore(
					substringAfter(originalUrl, "https://www.notion.so/image/"),
					"?table=",
				),
			)
			downloadImage(client, decoded, folder, flashcardWithOldImage.Id)
		} else {
			downloadImage(client, originalUrl, folder, flashcardWithOldImage.Id)
		}
	}
}

func substringAfter(s, sep string) string {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return ""
	}
	return s[idx+len(sep):]
}

func substringBefore(s, sep string) string {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func updateImagesInDb(db sqlx.DB, flashcard []flashcard.Flashcard, tableName string, folder string) {
	for _, flashcardWithOldImage := range flashcard {
		fileName, error := findFileByNameWithoutExt(folder, flashcardWithOldImage.Id)
		if error != nil {
			log.Println("ORIGINAL: " + *flashcardWithOldImage.Image)
			log.Println(error)
		}
		flashcardWithOldImage.UpdateImage(db, tableName, fileName)
	}
}
