package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"image_dock_server/handlers"
	"image_dock_server/internal"
	"image_dock_server/utils"
)

func main() {
	utils.LoadEnv()

	// Initialize database connection
	utils.InitDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s3Client := utils.NewS3Client()

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadHandler(w, r, s3Client)
	})
	mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListImagesHandler(w, r, s3Client)
	})

	handler := internal.EnableCORS(mux)

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
