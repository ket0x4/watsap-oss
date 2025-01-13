package messages

import "watsap/utils/telegram"

func StartupMessage1() {
	telegram.TgSendMsg(GetInitMsg())
}
