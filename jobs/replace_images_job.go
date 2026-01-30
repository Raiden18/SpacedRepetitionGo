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
	cache := image.LoadImageCache(imagesFolder)

	flashCardsToRevise := utils.Filter(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_REVISE_TABLE),
		flashcard.HasImage,
	)
	flashCardsToMemorize := utils.Filter(
		flashcard.NewFlashcardsFromDb(db, flashcard.FLASH_CARDS_TO_MEMORIZE_TABLE),
		flashcard.HasImage,
	)

	allFlashCards := append(flashCardsToMemorize, flashCardsToRevise...)
	keepBaseNames := map[string]struct{}{}
	for _, f := range allFlashCards {
		keepBaseNames[f.Id] = struct{}{}
	}
	image.DeleteFilesNotInSet(imagesFolder, keepBaseNames)
	for key := range cache.Entries {
		if _, ok := keepBaseNames[key]; !ok {
			delete(cache.Entries, key)
		}
	}

	downloadImages(client, imagesFolder, allFlashCards, &cache)

	firstRoundCoverters := map[string]func(path string){
		".jfif": image.ConvertJfifToJpg,
		".webp": image.ConvertWebpToJpg,
		".svg":  image.ConvertSvgToPng,
	}
	secondRoundCoverters := map[string]func(path string){
		".png": image.ConvertPngToJpg,
	}

	image.FindFilesAndConvert(imagesFolder, firstRoundCoverters)
	image.FindFilesAndConvert(imagesFolder, secondRoundCoverters)

	image.DoForEachFile(imagesFolder, image.ReduceImageSizeOfBigImage)

	updateCacheEntriesAfterConversion(&cache, imagesFolder)
	image.SaveImageCache(imagesFolder, cache)

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

const fallbackImageURL = "https://img.freepik.com/free-vector/hand-drawn-404-error_23-2147746234.jpg"

func downloadImage(client *resty.Client, url, folder, baseFilename string, cache *image.ImageCache) {
	entry, hasEntry := cache.Entries[baseFilename]
	if hasEntry && entry.URL != "" && entry.URL != url {
		hasEntry = false
		entry = image.ImageCacheEntry{}
		delete(cache.Entries, baseFilename)
	}

	fallbackUsed := false

	req := client.R().
		SetDoNotParseResponse(true).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36").
		SetHeader("Accept", "image/jpeg,image/webp,image/*,*/*")

	if hasEntry && entry.ETag != "" {
		req = req.SetHeader("If-None-Match", entry.ETag)
	}
	if hasEntry && entry.LastModified != "" {
		req = req.SetHeader("If-Modified-Since", entry.LastModified)
	}

	resp, err := req.
		Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() == 304 {
		needDownload := true
		if hasEntry && entry.FileName != "" {
			if _, err := os.Stat(filepath.Join(folder, entry.FileName)); err == nil {
				needDownload = false
			}
		} else if image.HasFileWithBaseName(folder, baseFilename) {
			needDownload = false
		}

		if !needDownload {
			return
		}

		resp.RawBody().Close()
		resp, err = client.R().
			SetDoNotParseResponse(true).
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36").
			SetHeader("Accept", "image/jpeg,image/webp,image/*,*/*").
			Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.RawBody().Close()
	}

	if resp.StatusCode() == 404 {
		resp.RawBody().Close()
		resp, err = client.R().
			SetDoNotParseResponse(true).
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36").
			SetHeader("Accept", "image/jpeg,image/webp,image/*,*/*").
			Get(fallbackImageURL)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.RawBody().Close()
		fallbackUsed = true
	}

	if !resp.IsSuccess() {
		log.Printf("image download failed with status %d for %s", resp.StatusCode(), url)
		return
	}

	ext := ".img"
	if contentType := resp.Header().Get("Content-Type"); contentType != "" {
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			ext = exts[0]
		}
	}

	fullPath := filepath.Join(folder, baseFilename+ext)

	tmpFile, err := os.CreateTemp(folder, baseFilename+".tmp-*")
	if err != nil {
		log.Println(err)
		return
	}
	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, resp.RawBody())
	if err != nil {
		log.Println(err)
		tmpFile.Close()
		os.Remove(tmpPath)
		return
	}

	if err := tmpFile.Close(); err != nil {
		log.Println(err)
		os.Remove(tmpPath)
		return
	}

	if err := os.Rename(tmpPath, fullPath); err != nil {
		log.Println(err)
		os.Remove(tmpPath)
		return
	}

	fileHash, err := image.FileSHA256(fullPath)
	if err != nil {
		log.Println(err)
	}

	cacheEntry := image.ImageCacheEntry{
		URL:      url,
		FileHash: fileHash,
		FileName: filepath.Base(fullPath),
	}
	if !fallbackUsed {
		cacheEntry.ETag = resp.Header().Get("ETag")
		cacheEntry.LastModified = resp.Header().Get("Last-Modified")
	}
	cache.Entries[baseFilename] = cacheEntry
}

func downloadImages(client *resty.Client, folder string, flashcard []flashcard.Flashcard, cache *image.ImageCache) {
	for _, flashcardWithOldImage := range flashcard {
		originalUrl := *flashcardWithOldImage.Image
		if strings.Contains(originalUrl, "base64") {
			sourceHash := image.BytesSHA256([]byte(originalUrl))
			if entry, ok := cache.Entries[flashcardWithOldImage.Id]; ok && entry.SourceHash == sourceHash {
				if image.HasFileWithBaseName(folder, flashcardWithOldImage.Id) {
					continue
				}
			}
			image.ConvertBase64ToImage(originalUrl, folder, flashcardWithOldImage.Id)
			if fileName, err := image.FindFileByNameWithoutExt(folder, flashcardWithOldImage.Id); err == nil {
				fileHash, err := image.FileSHA256(fileName)
				if err != nil {
					log.Println(err)
				}
				cache.Entries[flashcardWithOldImage.Id] = image.ImageCacheEntry{
					SourceHash: sourceHash,
					FileHash:   fileHash,
					FileName:   filepath.Base(fileName),
				}
			}
		} else if strings.HasPrefix(originalUrl, "https://www.notion.so/image/") {
			decoded, _ := url.QueryUnescape(
				utils.SubstringBefore(
					utils.SubstringAfter(originalUrl, "https://www.notion.so/image/"),
					"?table=",
				),
			)
			downloadImage(client, decoded, folder, flashcardWithOldImage.Id, cache)
		} else {
			downloadImage(client, originalUrl, folder, flashcardWithOldImage.Id, cache)
		}
	}
}

func updateImagesInDb(db sqlx.DB, flashcard flashcard.Flashcard, tableName string, folder string) {
	fileName, error := image.FindFileByNameWithoutExt(folder, flashcard.Id)
	if error != nil {
		log.Println("ORIGINAL: " + *flashcard.Image)
		log.Println(error)
		return
	}
	flashcard.UpdateImage(db, tableName, fileName)
}

func updateCacheEntriesAfterConversion(cache *image.ImageCache, folder string) {
	for baseName, entry := range cache.Entries {
		fileName, err := image.FindFileByNameWithoutExt(folder, baseName)
		if err != nil {
			continue
		}

		entry.FileName = filepath.Base(fileName)
		if fileHash, err := image.FileSHA256(fileName); err == nil {
			entry.FileHash = fileHash
		}

		cache.Entries[baseName] = entry
	}
}
