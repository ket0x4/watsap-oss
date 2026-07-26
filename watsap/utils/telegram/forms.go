package telegram

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"watsap/utils/logger"
)

// form data for file upload
func CreateForm(form map[string][]byte) (string, *bytes.Buffer, error) {
	logger.Debug("Telegram", "Starting form creation")
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)

	for key, val := range form {
		logger.Debug("Telegram", "Processing field: %s", key)
		if key == "document" && len(val) > 0 && val[0] == '@' {
			filePath := string(val[1:])
			logger.Debug("Telegram", "Uploading file: %s", filePath)
			file, err := os.Open(filePath)
			if err != nil {
				logger.Error("Telegram", "Error opening file %s: %v", filePath, err)
				mp.Close()
				return "", nil, err
			}
			part, err := mp.CreateFormFile(key, filepath.Base(filePath))
			if err != nil {
				file.Close()
				logger.Error("Telegram", "Error creating form file for %s: %v", filePath, err)
				mp.Close()
				return "", nil, err
			}
			_, err = io.Copy(part, file)
			file.Close()
			if err != nil {
				mp.Close()
				return "", nil, err
			}
		} else {
			part, err := mp.CreateFormField(key)
			if err != nil {
				mp.Close()
				return "", nil, err
			}
			part.Write(val)
			logger.Debug("Telegram", "Added field %s", key)
		}
	}
	mp.Close()
	logger.Debug("Telegram", "Form creation completed")
	return mp.FormDataContentType(), body, nil
}
