package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the multipart form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from the request body
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file from body", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get the current date
		today := time.Now()
		// Create a folder with today's date if it doesn't exist
		folderName := fmt.Sprintf("./files/%s", today.Format("2006-01-02"))
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			err := os.MkdirAll(folderName, 0755)
			if err != nil {
				http.Error(w, "Error creating directory", http.StatusInternalServerError)
				return
			}
		}

		// Create a file in the folder with the uploaded file's name
		fileName := fmt.Sprintf("%s/%s", folderName, handler.Filename)
		outFile, err := os.Create(fileName)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy the content from the uploaded file to the new file
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded and stored successfully")
	} else {
		// Handle invalid HTTP method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
