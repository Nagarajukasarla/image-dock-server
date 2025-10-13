package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	log.Println("Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables.")
	} else {
		log.Println(".env file loaded successfully")
	}

	// Log key environment variables (without sensitive data)
	if port := os.Getenv("PORT"); port != "" {
		log.Printf("PORT: %s", port)
	}
	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		log.Printf("S3_BUCKET: %s", bucket)
	}
	if endpoint := os.Getenv("S3_ENDPOINT"); endpoint != "" {
		log.Printf("S3_ENDPOINT: %s", endpoint)
	}
	log.Println("Environment setup completed")
}
