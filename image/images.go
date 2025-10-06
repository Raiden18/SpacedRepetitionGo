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

	err = imaging.Save(resized, filePath)
	if err != nil {
		log.Println("failed to save image: %v", err)
	}
}

func ConvertPngToJpg(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".png") + ".jpg"
	cmd := exec.Command("convert", filePath, newPath)
	error := cmd.Run()
	if error != nil {
		log.Println("Could not convert png to jpg", error)
	}
	os.Remove(filePath)
}

func ConvertJfifToJpg(filePath string) {
	newPath := strings.TrimSuffix(filePath, ".jfif") + ".jpg"
	err := os.Rename(filePath, newPath)
	if err != nil {
		log.Println("Could not convert Jfif to Jpg.", err)
	}
}

func ConvertSvgToPng(svgPath string) {
	ext := filepath.Ext(svgPath)
	base := strings.TrimSuffix(svgPath, ext)
	pngPath := base + ".png"

	cmd := exec.Command("rsvg-convert", "-f", "png", "-o", pngPath, svgPath)
	if err := cmd.Run(); err != nil {
		log.Println("Could not convert Svg to PNG", err)
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
