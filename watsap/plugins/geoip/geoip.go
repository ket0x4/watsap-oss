package geoip

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"watsap/utils/config"
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
		config.Logger("[GeoIP] Error getting IP info: "+err.Error(), "error")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger("[GeoIP] Error reading response body: "+err.Error(), "error")
		return
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		config.Logger("[GeoIP] Error unmarshalling JSON: "+err.Error(), "error")
		return
	}

	// save IP info to a JSON file
	file, err := os.Create(config.GeoFile)
	if err != nil {
		config.Logger("[GeoIP] Error creating JSON file: "+err.Error(), "error")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(ipInfo)
	if err != nil {
		config.Logger("[GeoIP] Error encoding JSON: "+err.Error(), "error")
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
