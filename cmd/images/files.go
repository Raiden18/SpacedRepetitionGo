package main

import (
	"fmt"
	"image/jpeg"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

func createFolderIfNotExist(folder string) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func clearFolder(folder string) {
	d, err := os.Open(folder)
	if err != nil {
		log.Println(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Println(err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(folder, name))
		if err != nil {
			log.Println(err)
		}
	}
}

func renameJfifToJpg(folder string) {
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".jfif") {
			newPath := strings.TrimSuffix(path, ".jfif") + ".jpg"
			err := os.Rename(path, newPath)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	})
}

func convertWebPtoJPEGInFolder(folder string) {
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
		}

		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".webp") {
			convertWebPtoJPEG(path)
		}
		return nil
	})
}

func convertWebPtoJPEG(inputPath string) {
	inFile, err := os.Open(inputPath)
	if err != nil {
		log.Println(err)
	}
	defer inFile.Close()

	img, err := webp.Decode(inFile)
	if err != nil {
		log.Println(err)
	}

	outputPath := strings.TrimSuffix(inputPath, ".webp") + ".jpg"
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Println(err)
	}
	defer outFile.Close()

	jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	if err := os.Remove(inputPath); err != nil {
		log.Println("delete .webp error: %w", err)
	}
}

func findFileByNameWithoutExt(folderPath, baseName string) (string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		nameWithoutExt := strings.TrimSuffix(name, ext)

		if nameWithoutExt == baseName {
			return filepath.Join(folderPath, name), nil
		}
	}

	return "", fmt.Errorf("file %q not found in %q", baseName, folderPath)
}
