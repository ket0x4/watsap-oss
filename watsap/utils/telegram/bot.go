package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"watsap/utils/config"
)

// var GODEBUG = "http1debug=1"
// var GODEBUG = "http2debug=1"

type TelegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Send message to Telegram
func TgSendMsg(msg string) {
	log.Println("Sending message...")
	
	payload := TelegramPayload{
		ChatID:    config.TG_CHAT_ID,
		Text:      msg,
		ParseMode: "HTML",
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling payload:", err)
		return
	}
	// Wipe the JSON payload from memory securely when done
	defer config.WipeMemory(jsonData)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	urlStr := fmt.Sprintf("%s%s/sendMessage", config.TgBotAPI, config.TG_BOT_TOKEN)
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending message:", err)
	} else {
		defer resp.Body.Close()
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
