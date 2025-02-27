package update

import (
	"fmt"
	"log"
	"watsap/utils/config"
)

// compare checks if the current version matches the remote version and downloads the update if necessary.
func compare() {
	log.Println("Downloading remote version information...")
	err := DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		log.Printf("Failed to download file: %v\n", err)
		return
	}

	// Assuming that UpdateParser() has been called earlier and the variables are set correctly.
	if currentVersion == remoteVersion {
		log.Println("Current version is up-to-date. No need to update.")
		return
	}

	fmt.Println("Current version is outdated. Downloading update...")
	err = DownloadFile(config.UpdateFile, config.UPDATE_URL)
	if err != nil {
		log.Printf("Failed to download update: %v\n", err)
		return
	}

	fmt.Println("Update downloaded successfully.")
}
