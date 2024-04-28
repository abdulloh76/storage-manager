package utils

import (
	"bytes"
	"io"
	"net/http"
)

// Function to combine files into one
func CombineFiles(fileURLs []string) (*bytes.Buffer, error) {
	var combinedBuffer bytes.Buffer

	for _, url := range fileURLs {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(&combinedBuffer, bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}

	return &combinedBuffer, nil
}
