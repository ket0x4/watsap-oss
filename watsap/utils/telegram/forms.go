package telegram

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

// form data for file upload
func CreateForm(form map[string]string) (string, io.Reader, error) {
	log.Println("Starting form creation")
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		log.Printf("Processing field: %s", key)
		if strings.HasPrefix(val, "@") {
			filePath := val[1:]
			log.Printf("Uploading file: %s", filePath)
			val = filePath
			file, err := os.Open(val)
			if err != nil {
				log.Printf("Error opening file %s: %v", val, err)
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				log.Printf("Error creating form file for %s: %v", val, err)
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
			log.Printf("Added field %s with value %s", key, val)
		}
	}
	log.Println("Form creation completed")
	return mp.FormDataContentType(), body, nil
}
