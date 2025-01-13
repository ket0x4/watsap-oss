package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"watsap/utils/config"
)

// var GODEBUG = "http1debug=1"
//var GODEBUG = "http2debug=1"

// Send message to Telegram
func TgSendMsg(msg string) {
	msg = url.QueryEscape(msg) // URL encoding
	config.Logger("Sending message: "+msg, "info")
	url := fmt.Sprintf("%s%s%s", config.TgBotAPI, config.TG_BOT_TOKEN, config.TgSendTextMsg+msg+"&parse_mode=HTML")
	if _, err := http.Get(url); err != nil {
		config.Logger("Error sending message: "+err.Error(), "error")
	}
	config.Logger("Message sent successfully", "info")
}

// Send file with caption to Telegram
func TgSendFile(filePath string, caption string) error {
	form := map[string]string{
		"document":   "@" + filePath,
		"chat_id":    config.TG_CHAT_ID,
		"caption":    caption,
		"parse_mode": "HTML", // markdown doesn't work well
	}

	config.Logger("Sending file: "+filePath+" with caption: "+caption, "info")

	ct, body, err := CreateForm(form)
	if err != nil {
		config.Logger("Error sending file: "+err.Error(), "error")
		return err
	}

	url := fmt.Sprintf("%s%s%s", config.TgBotAPI, config.TG_BOT_TOKEN, config.TgFileApiURL)
	_, err = http.Post(url, ct, body)
	if err != nil {
		config.Logger("Error sending file: "+err.Error(), "error")
		return err
	}

	config.Logger("File sent successfully", "info")
	return nil
}

/*
test bot
func TestBot() {
	fmt.Println("Testing bot...")
	TgSendFile("geo.w", "Test file")
}
*/
