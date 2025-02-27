package keylog

import "log"

func InitKeylog() {
	log.Println("Keylog plugin initialized")
	InitKeyboard()
}
