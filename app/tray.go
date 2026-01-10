package app

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TrayManager struct {
	ctx context.Context
}

func NewTrayManager(ctx context.Context) *TrayManager {
	return &TrayManager{ctx: ctx}
}

func (t *TrayManager) Setup() {
	log.Println("[Tray] System tray initialized - window will hide to tray on close")
}

func (t *TrayManager) ShowWindow() {
	runtime.WindowShow(t.ctx)
	log.Println("[Tray] Window shown")
}

func (t *TrayManager) HideWindow() {
	runtime.WindowHide(t.ctx)
	log.Println("[Tray] Window hidden")
}

func (t *TrayManager) ToggleWindow() {
	log.Println("[Tray] Toggling window visibility")
	t.ShowWindow()
}
