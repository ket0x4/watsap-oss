package geoip

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"watsap/utils/config"
	"watsap/utils/logger"
)

var (
	// variables
	UserIP      = ""
	UserCity    = ""
	UserRegion  = ""
	UserCountry = ""
	UserISP     = ""
	UserLoc     = ""
)

// get user external IP address and geo location
func GetIP() {
	type IPInfo struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Org      string `json:"org"`
		TimeZone string `json:"timezone"`
		Postal   string `json:"postal"`
		Loc      string `json:"loc"`
	}

	resp, err := http.Get("http://ipinfo.io/json")
	if err != nil {
		logger.Error("GeoIP", "Error getting IP info: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GeoIP", "Error reading response body: %s", err.Error())
		return
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		logger.Error("GeoIP", "Error unmarshalling JSON: %s", err.Error())
		return
	}

	// save IP info to a JSON file
	file, err := os.Create(config.GeoFile)
	if err != nil {
		logger.Error("GeoIP", "Error creating JSON file: %s", err.Error())
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(ipInfo)
	if err != nil {
		logger.Error("GeoIP", "Error encoding JSON: %s", err.Error())
		return
	}

	// set global variables
	UserIP = ipInfo.IP
	UserCity = ipInfo.City
	UserRegion = ipInfo.Region
	UserCountry = ipInfo.Country
	UserISP = ipInfo.Org
	UserLoc = ipInfo.Loc
}
