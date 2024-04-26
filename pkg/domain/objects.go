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

	var metadata types.MetadataModel = types.MetadataModel{
		FileExtension: ext,
		ObjectName:    fileHeader.Filename,
		MimeType:      mimeType,
		Size:          fileHeader.Size,
		StorageServer: f.storageServers[0], // todo
	}

	fileMetadataId, err := f.store.CreateMetadata(&metadata)
	if err != nil {
		return "", fmt.Errorf("error while saving metadata: %s", err)
	}

	fileNameFromMetadata := fmt.Sprintf("%s%s", fileMetadataId, ext)

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	writer.WriteField("from", "server-manager") // * just test key-value data
	// todo need to refactor transportation method
	fileWriter, err := writer.CreateFormFile("file", fileNameFromMetadata)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %s", err)
	}
	_, err = io.Copy(fileWriter, file)
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

	return fileNameFromMetadata, nil
}

func (f *Objects) GetObject(fileName string) (fileUrl string, fileMetadata *types.MetadataModel, err error) {
	uuidFileName, err := uuid.Parse(fileName)
	if err != nil {
		return "", nil, err
	}

	fileMetadata, err = f.store.GetMetadata(uuidFileName)
	if err != nil {
		return "", nil, err
	}

	fileUrl = f.storageServers[0] + f.downloadEndPoint + fileMetadata.ID.String() + fileMetadata.FileExtension // todo

	return fileUrl, fileMetadata, nil
}
