package handlers

import (
	"fmt"
	"net/http"

	"github.com/abdulloh76/storage-server/pkg/domain"
)

type HttpHandler struct {
	objects *domain.Objects
}

func NewHttpHandler(o *domain.Objects) *HttpHandler {

	return &HttpHandler{
		objects: o,
	}
}

func RegisterHandlers(handler *HttpHandler) {
	http.HandleFunc("/upload", handler.HandleUpload)
	http.HandleFunc("/download/", handler.HandleDownload)
}

func (h *HttpHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the multipart form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from the request body
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file from body", http.StatusBadRequest)
			return
		}
		defer file.Close()

		err = h.objects.UploadObject(file, fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded and stored successfully")
	} else {
		// Handle invalid HTTP method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HttpHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Extract the filename (without extension just the uuid) from the URL path
		fileName := r.URL.Path[len("/download/"):]
		if fileName == "" {
			http.Error(w, "File name not provided", http.StatusBadRequest)
			return
		}

		filePath, fileMetadata, err := h.objects.GetObject(fileName)

		w.Header().Set("Object-Name", fileMetadata.ObjectName)
		http.ServeFile(w, r, filePath)

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	} else {
		// Handle invalid HTTP method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
