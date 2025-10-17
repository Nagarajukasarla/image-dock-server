package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type ImageRecord struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename"`
	S3Key       string `json:"s3_key"`
	S3Bucket    string `json:"s3_bucket"`
	URL         string `json:"url"`
	UploadedAt  string `json:"uploaded_at"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	Name        string `json:"name"`
}

// InitDatabase initializes the database connection and creates tables if they don't exist
func InitDatabase() {
	log.Println(" Initializing database connection...")

	serviceURI := os.Getenv("DATABASE_URL")
	if serviceURI == "" {
		log.Fatal(" DATABASE_URL environment variable is required")
	}

	// Parse the database URL
	conn, err := url.Parse(serviceURI)
	if err != nil {
		log.Fatalf("❌ Failed to parse database URL: %v", err)
	}

	// Add SSL configuration for secure connection
	conn.RawQuery = "sslmode=require"

	// Initialize DB connection
	DB, err = sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatalf(" Failed to open database connection: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf(" Failed to ping database: %v", err)
	}

	log.Println(" Database connected successfully")

	// Create tables if they don't exist
	if err := createImagesTable(DB); err != nil {
		log.Fatalf(" Failed to create images table: %v", err)
	}

	log.Println(" Database initialization completed")
}

func createImagesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS images (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL,
		s3_key VARCHAR(500) NOT NULL,
		s3_bucket VARCHAR(255) NOT NULL,
		url TEXT,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		category VARCHAR(200) NOT NULL,
		sub_category VARCHAR(200) NOT NULL,
		name VARCHAR(255) NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create images table: %v", err)
	}

	log.Println("✅ Images table created/verified successfully")
	return nil
}

func StoreImage(filename, s3Key, s3Bucket, imageURL, category, subCategory, name string) error {
	query := `
		INSERT INTO images (filename, s3_key, s3_bucket, url, category, sub_category, name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := DB.Exec(query, filename, s3Key, s3Bucket, imageURL, category, subCategory, name)
	if err != nil {
		return fmt.Errorf("failed to store image record: %v", err)
	}

	log.Printf("💾 Image record stored in database - Filename: %s, S3Key: %s, Category: %s, SubCategory: %s, Name: %s", filename, s3Key, category, subCategory, name)
	return nil
}

func GetImageByID(id int) (*ImageRecord, error) {
	query := `SELECT id, filename, s3_key, s3_bucket, url, uploaded_at, category, sub_category, name FROM images WHERE id = $1`

	row := DB.QueryRow(query, id)

	var img ImageRecord
	err := row.Scan(&img.ID, &img.Filename, &img.S3Key, &img.S3Bucket, &img.URL, &img.UploadedAt, &img.Category, &img.SubCategory, &img.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	return &img, nil
}

func GetAllImages() ([]ImageRecord, error) {
	query := `SELECT id, filename, s3_key, s3_bucket, url, uploaded_at, category, sub_category, name FROM images ORDER BY uploaded_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %v", err)
	}
	defer rows.Close()

	var images []ImageRecord
	for rows.Next() {
		var img ImageRecord
		err := rows.Scan(&img.ID, &img.Filename, &img.S3Key, &img.S3Bucket, &img.URL, &img.UploadedAt, &img.Category, &img.SubCategory, &img.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image row: %v", err)
		}
		images = append(images, img)
	}

	return images, nil
}
