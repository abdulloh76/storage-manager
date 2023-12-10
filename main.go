package main

import (
	"fmt"
	"net/http"

	"github.com/abdulloh76/storage-server/pkg/domain"
	"github.com/abdulloh76/storage-server/pkg/handlers"
	"github.com/abdulloh76/storage-server/pkg/store"
)

func main() {
	DATABASE_URL := "postgres://user:password@localhost:5432/storage"
	postgresMetadataStore := store.NewPostgresDBStore(DATABASE_URL)

	objectDomain := domain.NewObjectsDomain(postgresMetadataStore)

	HttpHandler := handlers.NewHttpHandler(objectDomain)

	handlers.RegisterHandlers(HttpHandler)

	// Start the server on port 8080
	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
