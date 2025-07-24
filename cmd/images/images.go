package main

import (
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

func convertWebPtoJPEG(filePath string) {
	log.Println("HUI START")
	log.Println("FILE PATH: " + filePath)
	inFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer inFile.Close()

	img, err := webp.Decode(inFile)
	if err != nil {
		log.Println(err)
	}

	outputPath := strings.TrimSuffix(filePath, ".webp") + ".jpg"
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Println("Could not create out .jpg file", err)
	}
	defer outFile.Close()

	error := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	if error != nil {
		log.Println("Could not encode image.", error)
	}
	if err := os.Remove(filePath); err != nil {
		log.Println("delete .webp error: %w", err)
	}
	log.Println("HUI END")
}

func convertJfifToJpg(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".jfif") + ".jpg"
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Println(err)
	}
}

func convertSVGtoPNG(svgPath string) {
	ext := filepath.Ext(svgPath)
	base := strings.TrimSuffix(svgPath, ext)
	pngPath := base + ".png"

	cmd := exec.Command("rsvg-convert", "-f", "png", "-o", pngPath, svgPath)
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	os.Remove(svgPath)
}
