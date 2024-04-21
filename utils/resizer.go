package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

const pictureDirectory = "pictures"

func ResizeImage(image image.Image) image.Image {
	resizedImage := imaging.Resize(image, 320, 400, imaging.Linear)

	return resizedImage
}

func SaveLocal(image image.Image) (string, error) {
	timestamp := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s.jpg", timestamp)
	fullPath := filepath.Join(pictureDirectory, filename)

	err := os.MkdirAll(pictureDirectory, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory:", err)
		return "", err
	}

	imageFile, err := os.Create(fullPath)
	if err != nil {
		log.Println("Error creating output file:", err)
		return "", err
	}
	defer imageFile.Close()

	err = jpeg.Encode(imageFile, image, nil)
	if err != nil {
		log.Println("Error encoding image:", err)
		return "", err
	}

	return filename, err
}
