//go:build !windows
// +build !windows

package config

import (
	"fmt"
	"os"
)

// variables
var (
	PcName   = os.Getenv("HOSTNAME")
	UserName = os.Getenv("USER")
	homeDir  = getHomeDirUnix()
)

func getHomeDirUnix() string {
	home, err := os.UserHomeDir()

	if err != nil {
		home = os.Getenv("HOME")
	}

	return home
}

/* parse /etc/os-release
func getPlatform() string {
	platform := "unknown"
	osRelease := "/etc/os-release"
	if _, err := os.Stat(osRelease); err == nil {
		platform = "linux"
	}
	return platform
}
*/

// files & dirs
var waDirPrefix = fmt.Sprintf("%s/.config", homeDir)
