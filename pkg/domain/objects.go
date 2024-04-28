package domain

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/abdulloh76/storage-manager/pkg/types"
	"github.com/abdulloh76/storage-manager/pkg/utils"
	"github.com/google/uuid"
)

type Objects struct {
	store            types.MetadataStore
	storageServers   []string
	uploadEndPoint   string
	downloadEndPoint string
}

func NewObjectsDomain(s types.MetadataStore, storageServers []string, UPLOAD_ENDPOINT, DOWNLOAD_ENDPOINT string) *Objects {
	return &Objects{
		store:            s,
		storageServers:   storageServers,
		uploadEndPoint:   UPLOAD_ENDPOINT,
		downloadEndPoint: DOWNLOAD_ENDPOINT,
	}
}

func (f *Objects) UploadObject(file multipart.File, fileHeader *multipart.FileHeader) (filename string, err error) {
	ext := filepath.Ext(fileHeader.Filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "", fmt.Errorf("could not determine MIME type for extension %s", ext)
	}

	fileNameID := uuid.New()
	fileNameFromMetadata := fmt.Sprintf("%s%s", fileNameID, ext)

	buffers, err := utils.SplitFile(file, fileHeader, len(f.storageServers))
	for i := 0; i < len(buffers); i++ {
		var metadata types.MetadataModel = types.MetadataModel{
			FileNameID:    fileNameID,
			FileExtension: ext,
			ObjectName:    fileHeader.Filename,
			MimeType:      mimeType,
			Size:          int64(buffers[i].Len()),
			StorageServer: f.storageServers[i],
		}

		fileMetadataId, err := f.store.CreateMetadata(&metadata)
		if err != nil {
			return "", fmt.Errorf("error while saving metadata: %s", err)
		}

		requestBody := &bytes.Buffer{}
		writer := multipart.NewWriter(requestBody)
		writer.WriteField("from", "server-manager") // * just test key-value data
		// todo need to refactor transportation method
		fileWriter, err := writer.CreateFormFile("file", fileMetadataId.String()+".bin") // todo magic string
		if err != nil {
			return "", fmt.Errorf("error creating form file: %s", err)
		}
		_, err = io.Copy(fileWriter, buffers[i])
		if err != nil {
			return "", fmt.Errorf("error copying file contents: %s", err)
		}
		writer.Close()

		targetURL := metadata.StorageServer + f.uploadEndPoint
		req, err := http.NewRequest("POST", targetURL, requestBody)
		if err != nil {
			return "", fmt.Errorf("error creating request: %s", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		response, err := client.Do(req) // send the request
		if err != nil {
			return "", fmt.Errorf("error sending request: %s", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return "", fmt.Errorf("error from server %s: %s", metadata.StorageServer, err)
		}
	}

	return fileNameFromMetadata, nil
}

func (f *Objects) GetObject(fileName string) (file *bytes.Buffer, objectName string, err error) {
	uuidFileName, err := uuid.Parse(fileName)
	if err != nil {
		return nil, "", err
	}

	fileParts, err := f.store.GetMetadata(uuidFileName)
	if err != nil {
		return nil, "", err
	}

	filesURLs := make([]string, len(f.storageServers))
	for i := 0; i < len(fileParts); i++ {
		filesURLs[i] = fileParts[i].StorageServer + f.downloadEndPoint + fileParts[i].ID.String() + ".bin" // todo magic string
	}

	// todo everywhere we are trying to get [n] element of array even if we are not sure need to fix this
	objectName = fileParts[0].ObjectName

	file, err = utils.CombineFiles(filesURLs)
	if err != nil {
		return nil, "", err
	}

	return file, objectName, nil
}
