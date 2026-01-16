package app

import (
	"context"
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed tray-color.ico
var trayIconColor []byte

//go:embed tray-gray.ico
var trayIconGray []byte

type TrayManager struct {
	ctx            context.Context
	hasSound       bool
	lastUpdateTime time.Time
}

func NewTrayManager(ctx context.Context) *TrayManager {
	return &TrayManager{ctx: ctx}
}

func (t *TrayManager) Setup() {
	// Register systray without blocking main thread
	systray.Register(t.onReady, t.onExit)
	log.Println("[Tray] System tray registered")
}

func (t *TrayManager) SetHWND(hwnd uintptr) {
	// Not needed with getlantern/systray
}

func (t *TrayManager) Remove() {
	systray.Quit()
}

func (t *TrayManager) onReady() {
	// Set icon (gray by default - no sound)
	systray.SetIcon(trayIconGray)
	systray.SetTooltip("Round Sound Widget")
	t.hasSound = false

	// Add menu items - только выход
	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")

	// Handle menu clicks
	go func() {
		<-mQuit.ClickedCh
		log.Println("[Tray] Quit requested from tray menu")
		systray.Quit()
		os.Exit(0)
	}()

	log.Println("[Tray] System tray ready")
}

func (t *TrayManager) onExit() {
	log.Println("[Tray] System tray exited")
}

func (t *TrayManager) ShowWindow() {
	runtime.WindowShow(t.ctx)
	runtime.WindowUnminimise(t.ctx)
	log.Println("[Tray] Window shown")
}

func (t *TrayManager) HideWindow() {
	runtime.WindowHide(t.ctx)
	log.Println("[Tray] Window hidden")
}

func (t *TrayManager) ToggleWindow() {
	t.ShowWindow()
}

// SetIconState updates the tray icon based on sound presence (throttled)
func (t *TrayManager) SetIconState(hasSound bool) {
	// Throttle updates - no more than twice per second (500ms)
	now := time.Now()
	if now.Sub(t.lastUpdateTime) < 500*time.Millisecond {
		return
	}

	if t.hasSound == hasSound {
		return // No change needed
	}

	t.hasSound = hasSound
	t.lastUpdateTime = now

	if hasSound {
		systray.SetIcon(trayIconColor)
	} else {
		systray.SetIcon(trayIconGray)
	}
}
