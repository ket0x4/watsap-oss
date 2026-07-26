package wainit

import (
	"fmt"
	"os"
	"time"
	"watsap/utils/config"
	"watsap/utils/files"
	"watsap/utils/logger"

	"golang.org/x/exp/rand"
)

// assign user id
func AssignUserID() {
	rand.Seed(uint64(time.Now().UnixNano()))
	tempUid := fmt.Sprintf("%06d", rand.Intn(1000000))
	config.UserID = &tempUid
	// create user id file
	file, err := os.Create(config.DataFile)
	if err != nil {
		logger.Error("Init", "Error creating user id file: %s", err.Error())
		return
	}
	defer file.Close()
	// write user id to file
	_, err = file.WriteString(*config.UserID)
	if err != nil {
		logger.Error("Init", "Error writing user id to file: %s", err.Error())
		return
	}
}

// get user id from file
func GetUserID() string {
	if !files.Exists(config.DataFile) {
		return ""
	}
	// open user id file
	file, err := os.Open(config.DataFile)
	if err != nil {
		logger.Error("Init", "Error opening user id file: %s", err.Error())
		return ""
	}
	defer file.Close()
	// read user id from file
	buf := make([]byte, 6)
	_, err = file.Read(buf)
	if err != nil {
		logger.Error("Init", "Error reading user id from file: %s", err.Error())
		return ""
	}
	return string(buf)
}

// give finnal user id
func InitUserID() {
	if config.UserID == nil {
		tempUid := GetUserID()
		config.UserID = &tempUid
		if config.UserID == nil || *config.UserID == "" {
			AssignUserID()
		}
	}
	logger.Info("Init", "User ID: %s", *config.UserID)
}
