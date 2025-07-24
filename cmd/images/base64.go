package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func saveBase64ImageAutoExt(base64Str, outputFolder, baseFilename string) (string, error) {
	// Remove data URI prefix if present
	if idx := strings.Index(base64Str, "base64,"); idx != -1 {
		base64Str = base64Str[idx+7:]
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", errors.New("failed to decode base64: " + err.Error())
	}

	img, format, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return "", errors.New("failed to decode image: " + err.Error())
	}

	var ext string
	switch format {
	case "jpeg":
		ext = ".jpg"
	case "png":
		ext = ".png"
	case "gif":
		ext = ".gif"
	default:
		return "", errors.New("unsupported image format: " + format)
	}

	outputPath := outputFolder + "/" + baseFilename + ext

	// Save file
	switch format {
	case "jpeg":
		err = saveJPEG(outputPath, img)
	case "png":
		err = savePNG(outputPath, img)
	case "gif":
		err = os.WriteFile(outputPath, decoded, 0644)
	}

	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func saveJPEG(path string, img image.Image) error {
	f, err := os.CreateTemp("", "tmpjpeg")
	if err != nil {
		return err
	}
	defer f.Close()

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, img, nil)
}

func savePNG(path string, img image.Image) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}
