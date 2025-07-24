package main

import (
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
