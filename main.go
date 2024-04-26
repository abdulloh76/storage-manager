package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/abdulloh76/storage-manager/pkg/domain"
	"github.com/abdulloh76/storage-manager/pkg/handlers"
	"github.com/abdulloh76/storage-manager/pkg/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	DATABASE_URL := os.Getenv("DATABASE_URL")
	STORAGE_SERVERS := strings.Split(os.Getenv("STORAGE_SERVERS"), ",")
	UPLOAD_ENDPOINT := os.Getenv("UPLOAD_ENDPOINT")
	DOWNLOAD_ENDPOINT := os.Getenv("DOWNLOAD_ENDPOINT")

	postgresMetadataStore := store.NewPostgresDBStore(DATABASE_URL)
	objectDomain := domain.NewObjectsDomain(postgresMetadataStore, STORAGE_SERVERS, UPLOAD_ENDPOINT, DOWNLOAD_ENDPOINT)
	HttpHandler := handlers.NewHttpHandler(objectDomain)

	handlers.RegisterHandlers(HttpHandler)

	fmt.Println("Server running on http://localhost:" + PORT)
	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
