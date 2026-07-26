package update

import (
	"fmt"
	"log"
	"os"
	"watsap/utils/config"
)

func applyUpdate(newBinaryPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	backupPath := exePath + ".bak"
	_ = os.Remove(backupPath) // clean up any old backup

	err = os.Rename(exePath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to rename current executable to backup: %v", err)
	}

	err = os.Rename(newBinaryPath, exePath)
	if err != nil {
		_ = os.Rename(backupPath, exePath) // rollback
		return fmt.Errorf("failed to move new binary into place: %v", err)
	}

	err = os.Chmod(exePath, 0755)
	if err != nil {
		log.Printf("Warning: failed to set executable permissions: %v", err)
	}

	return nil
}

// compare checks if the current version matches the remote version and downloads the update if necessary.
func compare() {
	log.Println("Downloading remote version information...")
	err := DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		log.Printf("Failed to download remote version file: %v\n", err)
		return
	}

	_, err = UpdateParser()
	if err != nil {
		log.Printf("Failed to parse remote version info: %v\n", err)
		return
	}

	if currentVersion == remoteVersion {
		log.Println("Current version is up-to-date. No need to update.")
		return
	}

	log.Printf("Current version (%s) is outdated. New version (%s) is available. Downloading update...", currentVersion, remoteVersion)
	
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Failed to locate executable path: %v\n", err)
		return
	}
	newBinaryPath := exePath + ".new"

	err = DownloadFile(newBinaryPath, newURL)
	if err != nil {
		log.Printf("Failed to download update binary from %s: %v\n", newURL, err)
		return
	}

	log.Println("Applying update (hot-swapping binary)...")
	err = applyUpdate(newBinaryPath)
	if err != nil {
		log.Printf("Failed to apply update: %v\n", err)
		return
	}

	log.Println("Update applied successfully. Running version is updated.")
}

