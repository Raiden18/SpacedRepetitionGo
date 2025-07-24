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
	exec.Command("wkhtmltoimage", filePath, newPath).Run()
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
