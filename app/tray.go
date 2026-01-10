package app

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed tray.ico
var trayIconData []byte

type TrayManager struct {
	ctx context.Context
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
	// Set icon
	systray.SetIcon(trayIconData)
	systray.SetTooltip("Round Sound Widget")

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
