package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Extract the filename from the URL path
		fileName := r.URL.Path[len("/download/"):]
		if fileName == "" {
			http.Error(w, "File name not provided", http.StatusBadRequest)
			return
		}

		// Search for the file in all folders within the 'files' directory
		err := filepath.Walk("./files/", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() == fileName {
				http.ServeFile(w, r, path)
				return nil
			}
			return nil
		})

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	} else {
		// Handle invalid HTTP method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
