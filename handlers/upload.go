package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"image_dock_server/utils"
)

type UploadResponse struct {
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request, s3Client *s3.Client) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	bucket := os.Getenv("S3_BUCKET")
	uploadDir := os.Getenv("UPLOAD_DIR")

	// Sanitize filename by replacing spaces with underscores for S3 compatibility
	sanitizedFilename := strings.ReplaceAll(header.Filename, " ", "_")
	// Convert Windows path separators to forward slashes for S3
	key := strings.ReplaceAll(filepath.Join(uploadDir, sanitizedFilename), "\\", "/")

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("S3 upload failed: %v\n", err)
		http.Error(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	// Generate permanent public URL
	publicBase := os.Getenv("PUBLIC_URL_BASE")
	if publicBase == "" {
		log.Println("PUBLIC_URL_BASE not configured")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	publicURL := publicBase + "/" + strings.ReplaceAll(key, "\\", "/")

	// Get additional fields from form
	category := r.FormValue("category")
	subCategory := r.FormValue("sub_category")
	productName := r.FormValue("product_name")
	name := r.FormValue("name")

	// Store image record in database
	err = utils.StoreImage(header.Filename, key, bucket, publicURL, category, subCategory, productName, name)
	if err != nil {
		log.Printf("Failed to store image in database: %v", err)
		// Don't return error as upload was successful, just log it
	}

	resp := UploadResponse{
		Message: "Upload successful",
		URL:     publicURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ListImagesHandler(w http.ResponseWriter, r *http.Request, s3Client *s3.Client) {
	bucket := os.Getenv("S3_BUCKET")
	uploadDir := os.Getenv("UPLOAD_DIR")
	// Convert Windows path separators to forward slashes for S3
	prefix := strings.ReplaceAll(uploadDir, "\\", "/")

	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		http.Error(w, "Failed to list images", http.StatusInternalServerError)
		return
	}

	var urls []string
	publicBase := os.Getenv("PUBLIC_URL_BASE")
	if publicBase == "" {
		log.Println("PUBLIC_URL_BASE not configured")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	for _, item := range resp.Contents {
		publicURL := publicBase + "/" + strings.ReplaceAll(*item.Key, "\\", "/")
		urls = append(urls, publicURL)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
