package wainit

import (
	"fmt"
	"log"
	"os"
	"time"
	"watsap/utils/config"

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
		log.Printf("Error creating user id file: %s", err.Error())
		return
	}
	defer file.Close()
	// write user id to file
	_, err = file.WriteString(*config.UserID)
	if err != nil {
		log.Printf("Error writing user id to file: %s", err.Error())
		return
	}
}

// get user id from file
func GetUserID() string {
	// open user id file
	file, err := os.Open(config.DataFile)
	if err != nil {
		log.Printf("Error opening user id file: %s", err.Error())
		return ""
	}
	defer file.Close()
	// read user id from file
	buf := make([]byte, 6)
	_, err = file.Read(buf)
	if err != nil {
		log.Printf("Error reading user id from file: %s", err.Error())
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
	log.Printf("[Init] User ID: %s", *config.UserID)
}
