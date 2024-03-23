package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/abdulloh76/storage-manager/pkg/domain"
	"github.com/abdulloh76/storage-manager/pkg/handlers"
	"github.com/abdulloh76/storage-manager/pkg/store"
)

func main() {
	portFlag := flag.Int("port", 8080, "listening port")
	flag.Parse()
	PORT := fmt.Sprintf(":%d", *portFlag)

	DATABASE_URL := "postgres://user:password@localhost:5432/storage"
	postgresMetadataStore := store.NewPostgresDBStore(DATABASE_URL)

	objectDomain := domain.NewObjectsDomain(postgresMetadataStore)

	HttpHandler := handlers.NewHttpHandler(objectDomain)

	handlers.RegisterHandlers(HttpHandler)

	fmt.Println("Server running on http://localhost" + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
