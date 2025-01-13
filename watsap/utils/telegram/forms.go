package telegram

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"watsap/utils/config" // Import the logging package
)

// form data for file upload
func CreateForm(form map[string]string) (string, io.Reader, error) {
	config.Logger("Starting form creation", "info")
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		config.Logger(fmt.Sprintf("Processing field: %s", key), "info")
		if strings.HasPrefix(val, "@") {
			filePath := val[1:]
			config.Logger(fmt.Sprintf("Uploading file: %s", filePath), "info")
			val = filePath
			file, err := os.Open(val)
			if err != nil {
				config.Logger(fmt.Sprintf("Error opening file %s: %v", val, err), "error")
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				config.Logger(fmt.Sprintf("Error creating form file for %s: %v", val, err), "error")
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
			config.Logger(fmt.Sprintf("Added field %s with value %s", key, val), "info")
		}
	}
	config.Logger("Form creation completed", "info")
	return mp.FormDataContentType(), body, nil
}
