package files

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
)

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	log.Printf("Starting to copy from %s to %s", src, dst)
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Printf("Error opening source file: %v", err)
		return err
	}
	defer sourceFile.Close()

	// Get source file permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		log.Printf("Error retrieving source file info: %v", err)
		return err
	}

	// Create destination file with the same permissions
	destinationFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, sourceInfo.Mode())
	if err != nil {
		log.Printf("Error creating destination file: %v", err)
		return err
	}
	defer destinationFile.Close()

	buf := bufio.NewWriter(destinationFile)
	_, err = io.Copy(buf, sourceFile)
	if err != nil {
		log.Printf("Error copying data: %v", err)
		return err
	}
	if err := buf.Flush(); err != nil {
		log.Printf("Error flushing buffer: %v", err)
		return err
	}
	log.Printf("Successfully copied from %s to %s", src, dst)
	return nil
}

// copies files with the specified extensions from multiple srcDirs to dstDir
func copyFilesWithExtensions(srcDirs []string, dstDir string, extensions []string) error {
	for _, srcDir := range srcDirs {
		log.Printf("Copying files from directory: %s", srcDir)
		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				log.Printf("Error walking through %s: %v", srcDir, walkErr)
				return walkErr
			}

			if !info.IsDir() {
				for _, ext := range extensions {
					if filepath.Ext(info.Name()) == ext {
						relativePath, err := filepath.Rel(srcDir, path)
						if err != nil {
							log.Printf("Error getting relative path for %s: %v", path, err)
							return err
						}

						dstPath := filepath.Join(dstDir, relativePath)
						if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
							log.Printf("Error creating directories for %s: %v", dstPath, err)
							return err
						}

						if err := copyFile(path, dstPath); err != nil {
							log.Printf("Error copying file %s: %v", path, err)
							return err
						}
						break
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Failed to copy files from %s: %v", srcDir, err)
			return err
		}
	}
	return nil
}
