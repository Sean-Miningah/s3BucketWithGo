package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"fileUploadAWS/repo"
	"fileUploadAWS/utils"
)

func UploadHandler(config utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		file, _, err := r.FormFile("image")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error reading file: %v", err)
			return
		}
		defer file.Close()

		imageDoc, _, err := image.Decode(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding image: %v", err)
			return
		}

		resizedImage := utils.ResizeImage(imageDoc)
		filename, err := utils.SaveLocal(resizedImage)

		repo := repo.NewS3Client(
			config.AWS_S3_BUCKET_ACCESS_KEY,
			config.AWS_S3_BUCKET_SECRET_ACCESS_KEY,
			config.AWS_REGION,
		)

		// Create singed url used for uploading file
		presignedurl, err := repo.PutObject(config.AWS_BUCKET_NAME, filename, 60)
		if err != nil {
			log.Fatalf("Error generating presigned url for put object: %s", err)
		}
		// Uploading document to S3
		err = repo.UploadFile(resizedImage, presignedurl.URL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error uploading image to S3: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully uploaded image: %s", filename)
	}
}

func DeleteHandler(config utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		type RequestBody struct {
			Filename string `json:"filename"`
		}
		var data RequestBody
		if r.Body != nil {
			json.NewDecoder(r.Body).Decode(&data)
		}

		if data.Filename == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Filename must be provided")
			return
		}

		repo := repo.NewS3Client(
			config.AWS_S3_BUCKET_ACCESS_KEY,
			config.AWS_S3_BUCKET_SECRET_ACCESS_KEY,
			config.AWS_REGION,
		)

		presignedurl, err := repo.DeleteObject(config.AWS_BUCKET_NAME, data.Filename, 60)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error generating presigned url for put object: %s", err)
			return
		}

		err = repo.DeleteFile(presignedurl.URL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error deleting file from s3 bucket: %v", err)
			return
		}

		w.WriteHeader(200)
		fmt.Fprintf(w, "Successfully deleted file: %s", data.Filename)
	}
}

func setupRoutes() {
	// Load configurations from environment variables from .env
	// at root of project
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	config, err := utils.LoadViperEnvironment(cwd)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	http.HandleFunc("/upload", UploadHandler(config))
	http.HandleFunc("/delete", DeleteHandler(config))
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("hello world")
	setupRoutes()
}
