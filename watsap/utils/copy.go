package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// GetAutostartPath returns the autostart directory path based on the operating system.
func GetAutostartPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	case "linux":
		return filepath.Join(homeDir, ".config", "autostart")
	default:
		return ""
	}
}

// CopySelfToTempDir copies the current executable to the autostart directory and sets executable permissions.
func CopySelfToTempDir() error {
	log.Println("Starting CopySelfToTempDir")
	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Error getting executable path: %s", err.Error())
		return err
	}
	log.Println("Executable path:", exePath)

	// Create the destination path
	dstPath := filepath.Join(GetAutostartPath(), filepath.Base(exePath))
	log.Println("Destination path:", dstPath)

	// Open the source file
	srcFile, err := os.Open(exePath)
	if err != nil {
		log.Printf("Error opening source file: %s", err.Error())
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Error creating destination file: %s", err.Error())
		return err
	}
	defer dstFile.Close()

	// Copy the file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Printf("Error copying file contents: %s", err.Error())
		return err
	}

	// Set the executable permissions (only for Unix-like systems)
	if runtime.GOOS != "windows" {
		if err := dstFile.Chmod(0755); err != nil {
			log.Printf("Error setting executable permissions: %s", err.Error())
			return err
		}
	}

	log.Println("File copied to autostart directory successfully")
	return nil
}
