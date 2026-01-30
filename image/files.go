package image

import (
	"bytes"
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

func DeleteFilesNotInSet(folder string, keepBaseNames map[string]struct{}) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		log.Println(err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		if _, ok := keepBaseNames[base]; ok {
			continue
		}

		if err := os.Remove(filepath.Join(folder, name)); err != nil {
			log.Println(err)
		}
	}
}

func DoForEachFile(folder string, doFunc func(path string)) {
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		doFunc(path)
		return nil
	})
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

		if converter, ok := converters[".webp"]; ok {
			if IsWebpFile(path) {
				converter(path)
				return nil
			}
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

func HasFileWithBaseName(folderPath, baseName string) bool {
	_, err := FindFileByNameWithoutExt(folderPath, baseName)
	return err == nil
}

func IsWebpFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	header := make([]byte, 12)
	n, err := io.ReadFull(file, header)
	if err != nil || n < 12 {
		return false
	}

	return bytes.Equal(header[0:4], []byte("RIFF")) && bytes.Equal(header[8:12], []byte("WEBP"))
}

func saveImage(path string, image image.Image, encode func(w io.Writer, image image.Image) error) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if err := encode(tmpFile, image); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
