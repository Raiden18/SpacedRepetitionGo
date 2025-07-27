package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func convertHtmtoPng(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".htm") + ".png"
	cmd := exec.Command("wkhtmltoimage", filePath, newPath)
	error := cmd.Run()
	if error != nil {
		log.Println("Could not convert htm to png", error)
	}
	os.Remove(filePath)
}

func convertJfifToJpg(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".jfif") + ".jpg"
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Println("Could not convert Jfif to Jpg.", err)
	}
}

func convertSVGtoPNG(svgPath string) {
	ext := filepath.Ext(svgPath)
	base := strings.TrimSuffix(svgPath, ext)
	pngPath := base + ".png"

	cmd := exec.Command("rsvg-convert", "-f", "png", "-o", pngPath, svgPath)
	if err := cmd.Run(); err != nil {
		log.Println("Could not convert Svg to PNG", err)
	}

	os.Remove(svgPath)
}

func convertRngToPng(filePath string) {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	pngPath := base + ".png"

	cmd := exec.Command("convert", filePath, pngPath)
	if err := cmd.Run(); err != nil {
		log.Println("Could not convert RNG to PNG", err)
	}

	os.Remove(filePath)
}

func convertBase64ToImage(base64Str, outputFolder, baseFilename string) {
	if idx := strings.Index(base64Str, "base64,"); idx != -1 {
		base64Str = base64Str[idx+7:]
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Str)

	if err != nil {
		log.Println("Could not decode base64 image.", err)
		return
	}

	img, format, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		log.Println("Could not decode image.", err)
		return
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
		log.Println("Could not decode base64 image. Unsupported image formag: " + format)
	}

	outputPath := outputFolder + "/" + baseFilename + ext

	savingToFile := map[string]func(){
		"jpeg": func() { saveJPEG(outputPath, img) },
		"png":  func() { savePNG(outputPath, img) },
		"gif":  func() { saveGif(outputPath, decoded) },
	}

	save := savingToFile[format]
	save()
}

func savePNG(path string, img image.Image) {
	saveImage(path, img, png.Encode)
}

func saveGif(path string, decoded []byte) {
	os.WriteFile(path, decoded, 0644)
}

func saveJPEG(path string, img image.Image) {
	error := saveImage(
		path,
		img,
		func(w io.Writer, image image.Image) error {
			return jpeg.Encode(w, image, nil)
		},
	)
	if error != nil {
		log.Println("Could not save image in file. filePath : " + path)
	}
}
