package image

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

	"github.com/disintegration/imaging"
)

func ReduceImageSizeOfBigImage(filePath string) {
	img, err := imaging.Open(filePath)
	maxTelegramAllowedSize := 4096

	if err != nil {
		log.Println("failed to open image: "+filePath, err)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	var resized *image.NRGBA

	if width > maxTelegramAllowedSize || height > maxTelegramAllowedSize {
		resized = imaging.Fit(img, maxTelegramAllowedSize, maxTelegramAllowedSize, imaging.Lanczos)
	} else {
		resized = imaging.Clone(img)
	}

	dir := filepath.Dir(filePath)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(filePath)+".tmp-*")
	if err != nil {
		log.Println("failed to create temp image: "+filePath, err)
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	err = imaging.Save(resized, tmpPath)
	if err != nil {
		log.Printf("failed to save image: %v\n", err)
		os.Remove(tmpPath)
		return
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		log.Println("failed to replace image: "+filePath, err)
		os.Remove(tmpPath)
	}
}

func ConvertPngToJpg(filePath string) {
	ext := filepath.Ext(filePath)
	newPath := strings.TrimSuffix(filePath, ext) + ".jpg"
	tmpFile, err := os.CreateTemp(filepath.Dir(newPath), filepath.Base(newPath)+".tmp-*")
	if err != nil {
		log.Println("Could not create temp file for png conversion", err)
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	cmd := exec.Command("convert", filePath, tmpPath)
	error := cmd.Run()
	if error != nil {
		log.Println("Could not convert png to jpg", error)
		os.Remove(tmpPath)
		return
	}
	if err := os.Rename(tmpPath, newPath); err != nil {
		log.Println("Could not move converted jpg into place", err)
		os.Remove(tmpPath)
		return
	}
	os.Remove(filePath)
}

func ConvertWebpToJpg(filePath string) {
	ext := filepath.Ext(filePath)
	newPath := strings.TrimSuffix(filePath, ext) + ".jpg"
	if err := exec.Command("convert", filePath, newPath).Run(); err != nil {
		log.Println("Could not convert webp to jpg", err)
		return
	}
	os.Remove(filePath)
}

func ConvertJfifToJpg(filePath string) {
	ext := filepath.Ext(filePath)
	newPath := strings.TrimSuffix(filePath, ext) + ".jpg"
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Println("Could not convert Jfif to Jpg.", err)
	}
}

func ConvertSvgToPng(svgPath string) {
	ext := filepath.Ext(svgPath)
	base := strings.TrimSuffix(svgPath, ext)
	pngPath := base + ".png"

	tmpFile, err := os.CreateTemp(filepath.Dir(pngPath), filepath.Base(pngPath)+".tmp-*")
	if err != nil {
		log.Println("Could not create temp file for svg conversion", err)
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	cmd := exec.Command("rsvg-convert", "-f", "png", "-o", tmpPath, svgPath)
	if err := cmd.Run(); err != nil {
		log.Println("Could not convert Svg to PNG", err)
		os.Remove(tmpPath)
		return
	}

	if err := os.Rename(tmpPath, pngPath); err != nil {
		log.Println("Could not move converted png into place", err)
		os.Remove(tmpPath)
		return
	}

	os.Remove(svgPath)
}

func ConvertBase64ToImage(base64Str, outputFolder, baseFilename string) {
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

	extMap := map[string]string{
		"jpeg": ".jpg",
		"png":  ".png",
		"gif":  ".gif",
	}

	ext, hasKey := extMap[format]

	if !hasKey {
		log.Println("Could not decode base64 image. Unsupported image formag: " + format)
		return
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
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		log.Println("Could not create temp gif file", err)
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	if err := os.WriteFile(tmpPath, decoded, 0644); err != nil {
		log.Println("Could not save gif to temp file", err)
		os.Remove(tmpPath)
		return
	}

	if err := os.Rename(tmpPath, path); err != nil {
		log.Println("Could not move gif into place", err)
		os.Remove(tmpPath)
	}
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
