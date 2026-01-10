package app

import (
	"context"
	"log"
	"sync"
	"time"
)

// WindowManager handles window-related operations
type WindowManager struct {
	ctx       context.Context
	mu        sync.Mutex
	isRunning bool
	stopCh    chan struct{}
	Hwnd      uintptr // Window handle (Windows only)
}

// NewWindowManager creates a new WindowManager
func NewWindowManager(ctx context.Context) *WindowManager {
	return &WindowManager{
		ctx:    ctx,
		stopCh: make(chan struct{}),
	}
}

// SetDesktopLevel sets window to desktop level (behind all windows)
// This is a cross-platform stub, actual implementation is in window_windows.go
func (w *WindowManager) SetDesktopLevel() {
	w.setDesktopLevelImpl()
}

// StartDesktopLevelWatcher starts a goroutine to maintain desktop-level position
func (w *WindowManager) StartDesktopLevelWatcher() {
	w.mu.Lock()
	if w.isRunning {
		w.mu.Unlock()
		return
	}
	w.isRunning = true
	w.mu.Unlock()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	log.Println("Desktop level watcher started")

	for {
		select {
		case <-ticker.C:
			w.SetDesktopLevel()
		case <-w.stopCh:
			log.Println("Desktop level watcher stopped")
			return
		}
	}
}

// Stop stops the desktop level watcher
func (w *WindowManager) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.isRunning {
		close(w.stopCh)
		w.isRunning = false
	}
}
