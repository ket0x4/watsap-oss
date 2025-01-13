//go:build windows
// +build windows

package config

import (
	"os"
	"syscall"
)

// variables
var (
	PcName   = os.Getenv("COMPUTERNAME")
	UserName = os.Getenv("USERNAME")
)

func init() {
	if !DebugMode {
		Logger("[Wainit] Starting wainit", "default")
		const SW_HIDE = 0
		const SW_SHOW = 5
	}
}

// files & dirs
var waDirPrefix = os.Getenv("APPDATA")
var NTUserDir = os.Getenv("USERPROFILE")

// syscalls
var User32 = syscall.NewLazyDLL("user32.dll")
var ProcShowWindow = User32.NewProc("ShowWindow")
var ProcGetConsoleWindow = User32.NewProc("GetConsoleWindow")
