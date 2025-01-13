package files

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"watsap/utils/config"
)

func init() {
	dir := config.FilesDir
	err := compressDirectory(dir)
	if err != nil {
		config.Logger(fmt.Sprintf("[File] Error compressing files: %v", err), "error")
	} else {
		config.Logger("[File] Files compressed successfully", "default")
	}
}

func compressDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		if config.DebugMode {
			config.Logger("[File] No files to compress", "info")
		}
		return nil
	}

	archiveName := filepath.Join(config.WaDir, "files.zip")
	if _, err := os.Stat(archiveName); err == nil {
		config.Logger("[File] Archive already exists, skipping compression", "info")
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
			config.Logger(fmt.Sprintf("[File] Error adding %s to zip: %v", file.Name(), err), "error")
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
