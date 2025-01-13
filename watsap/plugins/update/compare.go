package update

import (
	"fmt"
	"watsap/utils/config"
)

// compare checks if the current version matches the remote version and downloads the update if necessary.
func compare() {
	fmt.Println("Downloading remote version information...")
	err := DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return
	}

	// Assuming that UpdateParser() has been called earlier and the variables are set correctly.
	if currentVersion == remoteVersion {
		fmt.Println("Current version is up-to-date. No need to update.")
		return
	}

	fmt.Println("Current version is outdated. Downloading update...")
	err = DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		fmt.Printf("Failed to download update: %v\n", err)
		return
	}

	fmt.Println("Update downloaded successfully.")
}
