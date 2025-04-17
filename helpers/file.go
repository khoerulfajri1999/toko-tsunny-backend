package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ImageKitSuccessResponse struct {
	URL string `json:"url"`
}

type ImageKitErrorResponse struct {
	Message string `json:"message"`
}

func UploadToImageKit(file io.Reader, fileName string) (string, error) {
	privateAPIKey := "private_Kair9oJ4QbD4qCtFVMaH0U+u0og="
	if privateAPIKey == "" {
		return "", errors.New("ImageKit private key not set")
	}

	endpoint := "https://upload.imagekit.io/api/v1/files/upload"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Field: file (actual file content)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	// Field: fileName (string)
	_ = writer.WriteField("fileName", fileName)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, &body)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(privateAPIKey, "")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errorResp ImageKitErrorResponse
		_ = json.Unmarshal(respBody, &errorResp)
		if errorResp.Message != "" {
			return "", errors.New("ImageKit Error: " + errorResp.Message)
		}
		return "", errors.New("ImageKit upload failed with status: " + resp.Status)
	}

	var successResp ImageKitSuccessResponse
	err = json.Unmarshal(respBody, &successResp)
	if err != nil {
		return "", err
	}

	return successResp.URL, nil
}
