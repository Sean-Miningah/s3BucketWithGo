package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"fileUploadAWS/repo"

	"github.com/disintegration/imaging"
)

const pictureDirectory = "pictures"
const MaxFileSize = 2 * 1024 * 1024 // 2MB

// ResizeAndSave resizes and saves the image to a PNG file
func ResizeAndSave(w http.ResponseWriter, r *http.Request) {
	// Open uploaded file
	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error reading file: %v", err)
		return
	}
	defer file.Close()

	imageFile, _, err := image.Decode(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding image: %v", err)
		return
	}

	finalImage := resizeToTargetSize(imageFile, 100)

	// A buffer to store the resized image
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, finalImage, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error endcoding image: %v", err)
		return
	}

	// Upload the resized image to s3
	repo := repo.NewS3Client()
	bucketName := "your_bucket_name"
	objectKey := "your_object_key"
	err = repo.UploadFile(&bucketName, &objectKey, "test_image.jpg", &buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error uploading image to S3: %v", err)
		return
	}

	// Save the resized image
	saveMessage := saveImage(finalImage)

	// Resize the image while preserving aspect ratio
	// resizedImage := imaging.Resize(imageFile, targetWidth, targetHeight, imaging.Lanczos)
	// saveMessage := saveImage(resizedImage)

	fmt.Fprintf(w, "Image resized and saved successfully! %s", saveMessage)
}

func resizeToTargetSize(img image.Image, targetSizeMB int) image.Image {
	// Adjust the dimensions until the size constraint is met (approximately 2MB)
	bit := 1024 * 1024
	targetSize := targetSizeMB * bit

	// for {
	// width := img.Bounds().Dx()
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	fmt.Printf("width: %d, height: %d\n", w, h)
	height := int(float64(w) * 3 / 4)
	// height := img.Bounds().Dy()
	width := w

	resizedImg := imaging.Thumbnail(img, width, height, imaging.Lanczos)
	// Check the size of the resized image
	beforeSize := imageSize(img) / bit
	size := imageSize(resizedImg) / bit
	log.Printf("before resizing : %v after resizing %v", beforeSize, size)
	if size <= targetSize {
		rate := size / beforeSize
		log.Printf("Image compression rate %v", rate)
	}
	return resizedImg
}

func imageSize(img image.Image) int {
	return img.Bounds().Dx() * img.Bounds().Dy() * 3 // Assuming 3 bytes per pixel (for RGB images)
}

func saveImage(image image.Image) string {
	// timestamp := time.Now().Format("20060102150405")
	// filename := fmt.Sprintf("%s.jpg", timestamp)
	filename := fmt.Sprintf("testImg-%v.jpg", 1)
	fullPath := filepath.Join(pictureDirectory, filename)

	err := os.MkdirAll(pictureDirectory, os.ModePerm) // Create directory if needed
	if err != nil {
		log.Println("Error creating directory:", err)
		return ""
	}

	outputFile, err := os.Create(fullPath)
	if err != nil {
		log.Println("Error creating output file:", err)
		return ""
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, image, nil)
	if err != nil {
		log.Println("Error encoding image:", err)
		return ""
	}

	return outputFile.Name()
}

func setupRoutes() {
	http.HandleFunc("/upload", ResizeAndSave)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("hello world")
	setupRoutes()
}
