package update

import (
	"fmt"
	"os"
	"watsap/utils/config"
	"watsap/utils/logger"
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
		logger.Warn("Update", "Failed to set executable permissions: %v", err)
	}

	return nil
}

// compare checks if the current version matches the remote version and downloads the update if necessary.
func compare() {
	if config.UPDATE_URL == "" {
		logger.Debug("Update", "UPDATE_URL not configured, skipping update check")
		return
	}
	logger.Info("Update", "Downloading remote version information...")
	err := DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		logger.Error("Update", "Failed to download remote version file: %v", err)
		return
	}

	_, err = UpdateParser()
	if err != nil {
		logger.Error("Update", "Failed to parse remote version info: %v", err)
		return
	}

	if currentVersion == remoteVersion {
		logger.Info("Update", "Current version is up-to-date (%s).", currentVersion)
		return
	}

	logger.Info("Update", "Current version (%s) is outdated. New version (%s) is available. Downloading update...", currentVersion, remoteVersion)
	
	exePath, err := os.Executable()
	if err != nil {
		logger.Error("Update", "Failed to locate executable path: %v", err)
		return
	}
	newBinaryPath := exePath + ".new"

	err = DownloadFile(newBinaryPath, newURL)
	if err != nil {
		logger.Error("Update", "Failed to download update binary from %s: %v", newURL, err)
		return
	}

	logger.Info("Update", "Applying update (hot-swapping binary)...")
	err = applyUpdate(newBinaryPath)
	if err != nil {
		logger.Error("Update", "Failed to apply update: %v", err)
		return
	}

	logger.Info("Update", "Update applied successfully. Running version is updated.")
}

