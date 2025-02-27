package files

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"watsap/utils/config"
)

func init() {
	dir := config.FilesDir
	err := compressDirectory(dir)
	if err != nil {
		log.Printf("[File] Error compressing files: %v", err)
	} else {
		log.Println("[File] Files compressed successfully")
	}
}

func compressDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		log.Println("[File] No files to compress")
		return nil
	}

	archiveName := filepath.Join(config.WaDir, "files.zip")
	if _, err := os.Stat(archiveName); err == nil {
		log.Println("[File] Archive already exists, skipping compression")
		return nil
	}

	zipFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if err := addFileToZip(zipWriter, filePath); err != nil {
			log.Printf("[File] Error adding %s to zip: %v", file.Name(), err)
			continue
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
