package telegram

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"watsap/utils/config"
)

// var GODEBUG = "http1debug=1"
// var GODEBUG = "http2debug=1"

// Send message to Telegram
func TgSendMsg(msg string) {
	msg = url.QueryEscape(msg) // URL encoding
	log.Println("Sending message:", msg)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	urlStr := fmt.Sprintf("%s%s%s", config.TgBotAPI, config.TG_BOT_TOKEN, config.GetTgSendTextMsg()+msg+"&parse_mode=HTML")
	if _, err := client.Get(urlStr); err != nil {
		log.Println("Error sending message:", err)
	} else {
		log.Println("Message sent successfully")
	}
}

// Send file with caption to Telegram
func TgSendFile(filePath string, caption string) error {
	form := map[string]string{
		"document":   "@" + filePath,
		"chat_id":    config.TG_CHAT_ID,
		"caption":    caption,
		"parse_mode": "HTML", // markdown doesn't work well
	}

	log.Println("Sending file:", filePath, "with caption:", caption)

	ct, body, err := CreateForm(form)
	if err != nil {
		log.Println("Error sending file:", err)
		return err
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Timeout slightly longer for files
	}

	urlStr := fmt.Sprintf("%s%s%s", config.TgBotAPI, config.TG_BOT_TOKEN, config.TgFileApiURL)
	_, err = client.Post(urlStr, ct, body)
	if err != nil {
		log.Println("Error sending file:", err)
		return err
	}

	log.Println("File sent successfully")
	return nil
}

/*
test bot
func TestBot() {
	fmt.Println("Testing bot...")
	TgSendFile("geo.w", "Test file")
}
*/
