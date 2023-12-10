package domain

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/abdulloh76/storage-server/pkg/types"
	"github.com/google/uuid"
)

type Objects struct {
	store types.MetadataStore
}

func NewObjectsDomain(s types.MetadataStore) *Objects {
	return &Objects{
		store: s,
	}
}

func (f *Objects) UploadObject(file multipart.File, fileHeader *multipart.FileHeader) error {
	// Get the current date
	today := time.Now()
	// Create a folder with today's date if it doesn't exist
	folderName := fmt.Sprintf("./files/%s", today.Format("2006-01-02"))

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.MkdirAll(folderName, 0755)
		if err != nil {
			return fmt.Errorf("Error creating directory")
		}
	}

	ext := filepath.Ext(fileHeader.Filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return fmt.Errorf("could not determine MIME type for extension %s", ext)
	}

	var metadata types.MetadataModel = types.MetadataModel{
		FileExtension: ext,
		ObjectName:    fileHeader.Filename,
		MimeType:      mimeType,
		Size:          fileHeader.Size,
	}

	fileMetadataId, err := f.store.CreateMetadata(&metadata)

	if err != nil {
		return fmt.Errorf("Error while saving metadata")
	}

	fileNameFromMetadata := fmt.Sprintf("%s%s", fileMetadataId, ext)

	// Create a file in the folder with the uploaded file's name
	fileName := fmt.Sprintf("%s/%s", folderName, fileNameFromMetadata)
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file")
	}
	defer outFile.Close()

	// Copy the content from the uploaded file to the new file
	_, err = io.Copy(outFile, file)
	if err != nil {
		return fmt.Errorf("Error copying file")
	}

	return nil
}

func (f *Objects) GetObject(fileName string) (filePath string, fileMetadata *types.MetadataModel, err error) {
	uuidFileName, err := uuid.Parse(fileName)
	if err != nil {
		return "", nil, err
	}

	fileMetadata, err = f.store.GetMetadata(uuidFileName)
	if err != nil {
		return "", nil, err
	}

	fileNameFromMetadata := fmt.Sprintf("%s%s", fileMetadata.ID, fileMetadata.FileExtension)

	// Search for the file in all folders within the 'files' directory
	err = filepath.Walk("./files/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileNameFromMetadata {
			filePath = path
			return nil
		}
		return nil
	})

	if err != nil {
		return "", nil, err
	}

	return filePath, fileMetadata, nil
}
