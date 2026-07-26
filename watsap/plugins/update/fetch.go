package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"watsap/utils/logger"
)

// DownloadFile downloads a file from the given URL and saves it to the specified path.
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Update", "Failed to download file from %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Update", "Bad status when downloading from %s: %s", url, resp.Status)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		logger.Error("Update", "Failed to create file %s: %v", filepath, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Error("Update", "Failed to copy content to %s: %v", filepath, err)
		return err
	}

	logger.Info("Update", "Successfully downloaded %s to %s", url, filepath)
	return nil
}
