package files

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"watsap/utils/config"
)

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	config.Logger(fmt.Sprintf("Starting to copy from %s to %s", src, dst), "info")
	sourceFile, err := os.Open(src)
	if (err != nil) {
		config.Logger(fmt.Sprintf("Error opening source file: %v", err), "error")
		return err
	}
	defer sourceFile.Close()

	// Get source file permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		config.Logger(fmt.Sprintf("Error retrieving source file info: %v", err), "error")
		return err
	}

	// Create destination file with the same permissions
	destinationFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, sourceInfo.Mode())
	if err != nil {
		config.Logger(fmt.Sprintf("Error creating destination file: %v", err), "error")
		return err
	}
	defer destinationFile.Close()

	buf := bufio.NewWriter(destinationFile)
	_, err = io.Copy(buf, sourceFile)
	if err != nil {
		config.Logger(fmt.Sprintf("Error copying data: %v", err), "error")
		return err
	}
	if err := buf.Flush(); err != nil {
		config.Logger(fmt.Sprintf("Error flushing buffer: %v", err), "error")
		return err
	}
	config.Logger(fmt.Sprintf("Successfully copied from %s to %s", src, dst), "info")
	return nil
}

// copies files with the specified extensions from multiple srcDirs to dstDir
func copyFilesWithExtensions(srcDirs []string, dstDir string, extensions []string) error {
	for _, srcDir := range srcDirs {
		config.Logger(fmt.Sprintf("Copying files from directory: %s", srcDir), "info")
		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				config.Logger(fmt.Sprintf("Error walking through %s: %v", srcDir, walkErr), "error")
				return walkErr
			}

			if !info.IsDir() {
				for _, ext := range extensions {
					if filepath.Ext(info.Name()) == ext {
						relativePath, err := filepath.Rel(srcDir, path)
						if err != nil {
							config.Logger(fmt.Sprintf("Error getting relative path for %s: %v", path, err), "error")
							return err
						}

						dstPath := filepath.Join(dstDir, relativePath)
						if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
							config.Logger(fmt.Sprintf("Error creating directories for %s: %v", dstPath, err), "error")
							return err
						}

						if err := copyFile(path, dstPath); err != nil {
							config.Logger(fmt.Sprintf("Error copying file %s: %v", path, err), "error")
							return err
						}
						break
					}
				}
			}
			return nil
		})
		if err != nil {
			config.Logger(fmt.Sprintf("Failed to copy files from %s: %v", srcDir, err), "error")
			return err
		}
	}
	return nil
}
