package utils

import (
	"bytes"
	"io"
	"mime/multipart"
)

// Function to split a file into n parts
func SplitFile(file multipart.File, fileHeader *multipart.FileHeader, numParts int) ([]*bytes.Buffer, error) {
	fileSize := fileHeader.Size
	chunkSize := (fileSize + int64(numParts) - 1) / int64(numParts)

	buffers := make([]*bytes.Buffer, numParts)
	for i := 0; i < numParts; i++ {
		buffer := new(bytes.Buffer)
		written, err := io.CopyN(buffer, file, chunkSize)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if written == 0 {
			break
		}
		buffers[i] = buffer
	}
	return buffers, nil
}
