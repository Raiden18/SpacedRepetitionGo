package main

import (
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
)

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
