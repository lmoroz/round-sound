//go:build windows

package app

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
	procFindWindowW   = user32.NewProc("FindWindowW")
	procGetWindowLong = user32.NewProc("GetWindowLongW")
	procSetWindowLong = user32.NewProc("SetWindowLongW")
	procEnumWindows   = user32.NewProc("EnumWindows")
	procGetWindowText = user32.NewProc("GetWindowTextW")
)

const (
	HWND_BOTTOM = 1

	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOACTIVATE     = 0x0010
	SWP_NOSENDCHANGING = 0x0400

	WS_EX_TOOLWINDOW = 0x00000080
	WS_EX_APPWINDOW  = 0x00040000
)

// GWL_EXSTYLE for 64-bit systems (negative value as unsigned)
const GWL_EXSTYLE = ^uintptr(19) // equivalent to -20

// cachedHWND stores the found window handle
var cachedHWND uintptr

// findRoundSoundWindow finds the Round Sound window handle
func findRoundSoundWindow() uintptr {
	if cachedHWND != 0 {
		return cachedHWND
	}

	var foundHWND uintptr

	// Callback function for EnumWindows
	callback := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		// Get window title
		title := make([]uint16, 256)
		procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&title[0])), 256)
		titleStr := syscall.UTF16ToString(title)

		if titleStr == "Round Sound" {
			foundHWND = hwnd
			return 0 // Stop enumeration
		}
		return 1 // Continue enumeration
	})

	procEnumWindows.Call(callback, 0)

	if foundHWND != 0 {
		cachedHWND = foundHWND
		log.Printf("Found Round Sound window: HWND=%d", foundHWND)
	}

	return foundHWND
}

// setDesktopLevelImpl sets window to desktop level using Windows API
func (w *WindowManager) setDesktopLevelImpl() {
	hwnd := findRoundSoundWindow()
	if hwnd == 0 {
		return
	}

	// Save hwnd for later use
	if w.Hwnd == 0 {
		w.Hwnd = hwnd
	}

	// Set window position to bottom of Z-order
	ret, _, err := procSetWindowPos.Call(
		hwnd,
		HWND_BOTTOM,
		0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE|SWP_NOSENDCHANGING,
	)

	if ret == 0 {
		log.Printf("SetWindowPos failed: %v", err)
	}
}

// setToolWindow makes window a tool window (no taskbar icon)
func setToolWindow(hwnd uintptr) {
	// Get current extended style
	exStyle, _, _ := procGetWindowLong.Call(hwnd, GWL_EXSTYLE)

	// Add WS_EX_TOOLWINDOW and remove WS_EX_APPWINDOW
	newStyle := (exStyle | WS_EX_TOOLWINDOW) &^ WS_EX_APPWINDOW

	procSetWindowLong.Call(hwnd, GWL_EXSTYLE, newStyle)
}
