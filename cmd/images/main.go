package main

import (
	"log"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/image"
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
	image.CreateFolderIfNotExist(imagesFolder)
	image.DeleteAllFilesFromFolder(imagesFolder)

	flashCardsToRevise := utils.Filter(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE),
		flashcard.HasImage,
	)
	flashCardsToMemorize := utils.Filter(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
		flashcard.HasImage,
	)

	allFlashCards := append(flashCardsToMemorize, flashCardsToRevise...)

	downloadImages(client, imagesFolder, allFlashCards)

	converts := map[string]func(path string){
		".jfif": image.ConvertJfifToJpg,
		".htm":  image.ConvertHtmtoPng,
		".svg":  image.ConvertSVGtoPNG,
	}

	image.FindFilesAndConvert(imagesFolder, converts)

	utils.ForEach(
		flashCardsToRevise,
		func(f flashcard.Flashcard) {
			updateImagesInDb(db, f, flashcard.FLASH_CARDS_TO_REVISE_TABLE, imagesFolder)
		},
	)

	utils.ForEach(
		flashCardsToRevise,
		func(f flashcard.Flashcard) {
			updateImagesInDb(db, f, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE, imagesFolder)
		},
	)
}

func downloadImages(client *resty.Client, folder string, flashcard []flashcard.Flashcard) {
	for _, flashcardWithOldImage := range flashcard {
		originalUrl := *flashcardWithOldImage.Image
		if strings.Contains(originalUrl, "base64") {
			image.ConvertBase64ToImage(originalUrl, folder, flashcardWithOldImage.Id)
		} else if strings.HasPrefix(originalUrl, "https://www.notion.so/image/") {
			decoded, _ := url.QueryUnescape(
				utils.SubstringBefore(
					utils.SubstringAfter(originalUrl, "https://www.notion.so/image/"),
					"?table=",
				),
			)
			downloadImage(client, decoded, folder, flashcardWithOldImage.Id)
		} else {
			downloadImage(client, originalUrl, folder, flashcardWithOldImage.Id)
		}
	}
}

func updateImagesInDb(db sqlx.DB, flashcard flashcard.Flashcard, tableName string, folder string) {
	fileName, error := image.FindFileByNameWithoutExt(folder, flashcard.Id)
	if error != nil {
		log.Println("ORIGINAL: " + *flashcard.Image)
		log.Println(error)
	}
	flashcard.UpdateImage(db, tableName, fileName)
}
