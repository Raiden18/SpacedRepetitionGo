package image

import (
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateFolderIfNotExist(folder string) {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func DeleteAllFilesFromFolder(folder string) {
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

func FindFilesAndConvert(folder string, converters map[string]func(path string)) {
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		lowerName := strings.ToLower(d.Name())

		for ext, converter := range converters {
			if strings.HasSuffix(lowerName, ext) {
				converter(path)
				break
			}
		}

		return nil
	})
}

func FindFileByNameWithoutExt(folderPath, baseName string) (string, error) {
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

func saveImage(path string, image image.Image, encode func(w io.Writer, image image.Image) error) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return encode(outFile, image)
}
