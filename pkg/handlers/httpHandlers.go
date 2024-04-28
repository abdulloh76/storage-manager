package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/abdulloh76/storage-manager/pkg/domain"
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
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file from body", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename, err := h.objects.UploadObject(file, fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded and stored successfully: %s", filename)
	} else {
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

		file, objectName, err := h.objects.GetObject(fileName)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.Header().Set("Object-Name", objectName)

		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Error copying file contents to client", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
