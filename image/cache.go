package image

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

const cacheFileName = ".image-cache.json"

type ImageCacheEntry struct {
	URL          string `json:"url,omitempty"`
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	SourceHash   string `json:"source_hash,omitempty"`
	FileHash     string `json:"file_hash,omitempty"`
	FileName     string `json:"file_name,omitempty"`
}

type ImageCache struct {
	Entries map[string]ImageCacheEntry `json:"entries"`
}

func LoadImageCache(folder string) ImageCache {
	cachePath := filepath.Join(folder, cacheFileName)
	cache := ImageCache{Entries: map[string]ImageCacheEntry{}}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cache
		}
		log.Println(err)
		return cache
	}

	if err := json.Unmarshal(data, &cache); err != nil {
		log.Println(err)
		return ImageCache{Entries: map[string]ImageCacheEntry{}}
	}

	if cache.Entries == nil {
		cache.Entries = map[string]ImageCacheEntry{}
	}

	return cache
}

func SaveImageCache(folder string, cache ImageCache) {
	cachePath := filepath.Join(folder, cacheFileName)
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	tmpPath := cachePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		log.Println(err)
		return
	}

	if err := os.Rename(tmpPath, cachePath); err != nil {
		log.Println(err)
	}
}

func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func BytesSHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
