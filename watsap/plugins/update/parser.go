package update

import (
	"encoding/json"
	"fmt"
	"os"
	"watsap/utils/config"
)

const jsonFile = "watsap.json"

var currentVersion string
var remoteVersion string
var sha256 string
var changeURL bool
var newURL string

// Update represents the structure of the JSON data.
type Update struct {
	CURRENT_VERSION string `json:"CURRENT_VERSION"`
	SHA256          string `json:"SHA256"`
	CHANGE_URL      bool   `json:"CHANGE_URL"`
	NEW_URL         string `json:"NEW_URL"`
}

// ParseUpdate parses JSON data into an Update struct.
func ParseUpdate(data []byte) (Update, error) {
	var u Update
	err := json.Unmarshal(data, &u)
	return u, err
}

// UpdateParser reads JSON data from a file, parses it, and sets new variables based on the parsed data.
func UpdateParser() (Update, error) {
	// Read JSON data from file
	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		return Update{}, fmt.Errorf("failed to read JSON file: %v", err)
	}

	// Parse the JSON data
	updateJson, err := ParseUpdate(jsonData)
	if err != nil {
		return Update{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Set new variables after parsing
	currentVersion = config.WaVersion
	remoteVersion = updateJson.CURRENT_VERSION
	sha256 = updateJson.SHA256
	changeURL = updateJson.CHANGE_URL
	newURL = updateJson.NEW_URL

	fmt.Printf("Current Version: %s\n", currentVersion)
	fmt.Printf("SHA256: %s\n", sha256)
	fmt.Printf("Change URL: %t\n", changeURL)
	fmt.Printf("New URL: %s\n", newURL)

	return updateJson, nil
}

func initUpdateParser() {
	_, err := UpdateParser()
	if err != nil {
		fmt.Printf("Error updating: %v\n", err)
	}
}
