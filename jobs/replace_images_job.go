package jobs

import (
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"spacedrepetitiongo/flashcard"
	"spacedrepetitiongo/image"
	"spacedrepetitiongo/utils"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
)

func ReplaceImages() {
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

	firstRoundCoverters := map[string]func(path string){
		".jfif": image.ConvertJfifToJpg,
		".htm":  image.ConvertHtmtoPng,
		".svg":  image.ConvertSvgToPng,
	}
	secondRoundCoverters := map[string]func(path string){
		".png": image.ConvertPngToJpg,
	}

	image.FindFilesAndConvert(imagesFolder, firstRoundCoverters)
	image.FindFilesAndConvert(imagesFolder, secondRoundCoverters)

	utils.ForEach(
		flashCardsToRevise,
		func(f flashcard.Flashcard) {
			updateImagesInDb(db, f, flashcard.FLASH_CARDS_TO_REVISE_TABLE, imagesFolder)
		},
	)

	utils.ForEach(
		flashCardsToMemorize,
		func(f flashcard.Flashcard) {
			updateImagesInDb(db, f, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE, imagesFolder)
		},
	)
}

func downloadImage(client *resty.Client, url, folder, baseFilename string) {
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36").
		SetHeader("Accept", "image/jpeg,image/webp,image/*,*/*").
		Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.RawBody().Close()

	ext := ".img"
	if contentType := resp.Header().Get("Content-Type"); contentType != "" {
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			ext = exts[0]
		}
	}

	fullPath := filepath.Join(folder, baseFilename+ext)

	out, err := os.Create(fullPath)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.RawBody())
	if err != nil {
		log.Println(err)
	}
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
