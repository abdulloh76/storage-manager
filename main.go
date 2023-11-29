package main

import (
	"fmt"
	"net/http"

	"github.com/abdulloh76/storage-server/handlers"
)

func main() {
	http.HandleFunc("/upload", handlers.HandleUpload)
	http.HandleFunc("/download/", handlers.HandleDownload)

	// Start the server on port 8080
	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
