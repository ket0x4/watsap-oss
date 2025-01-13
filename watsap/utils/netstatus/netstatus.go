package netstatus

import "net/http"

// Check connection to a given URL
func CheckConnection(url string) bool {
	_, err := http.Get(url)
	return err == nil
}

// TGbotapi connection check
func CheckTgApi() bool {
	return CheckConnection("https://api.telegram.org")
}

// Check network connection
func CheckNetwork() bool {
	return CheckConnection("https://www.google.com")
}

// check ipinfo.io connection
func CheckIpApi() bool {
	return CheckConnection("http://ipinfo.io")
}
