package update

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// DownloadFile downloads a file from the given URL and saves it to the specified path.
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download file from %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad status when downloading from %s: %s", url, resp.Status)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filepath, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Failed to copy content to %s: %v", filepath, err)
		return err
	}

	log.Printf("Successfully downloaded %s to %s", url, filepath)
	return nil
}
