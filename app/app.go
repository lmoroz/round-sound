package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"round-sound/media"
)

// App struct holds the application state
type App struct {
	ctx           context.Context
	mu            sync.RWMutex
	config        *Config
	wnpServer     *media.WebNowPlayingServer
	activePlayer  *media.Player
	windowManager *WindowManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg := LoadConfig()
	return &App{
		config: cfg,
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize window manager
	a.windowManager = NewWindowManager(ctx)

	// Start WebNowPlaying server
	var err error
	a.wnpServer, err = media.NewWebNowPlayingServer(8974, a.onPlayerUpdate)
	if err != nil {
		log.Printf("Failed to start WebNowPlaying server: %v", err)
		runtime.EventsEmit(ctx, "error:port_busy", "Порт 8974 занят другим приложением (возможно, Rainmeter WebNowPlaying). Закройте другое приложение и перезапустите Round Sound.")
	}

	// Start desktop-level window manager (HWND_BOTTOM)
	go a.windowManager.StartDesktopLevelWatcher()

	log.Println("Round Sound started")
}

// DomReady is called after the front-end dom has been loaded
func (a *App) DomReady(ctx context.Context) {
	// Set saved window position
	if a.config.WindowX != 0 || a.config.WindowY != 0 {
		runtime.WindowSetPosition(ctx, a.config.WindowX, a.config.WindowY)
		log.Printf("Restored window position: %d, %d", a.config.WindowX, a.config.WindowY)
	}

	// Set window to desktop level
	a.windowManager.SetDesktopLevel()

	// Start periodic position saver (every 5 seconds)
	go a.startPositionSaver()
}

// Shutdown is called when the app is closing
func (a *App) Shutdown(ctx context.Context) {
	// Save window position
	x, y := runtime.WindowGetPosition(ctx)
	a.config.WindowX = x
	a.config.WindowY = y
	a.config.Save()

	// Stop WebNowPlaying server
	if a.wnpServer != nil {
		a.wnpServer.Stop()
	}

	log.Println("Round Sound shutdown")
}

// onPlayerUpdate is called when player state changes
func (a *App) onPlayerUpdate(player *media.Player) {
	a.mu.Lock()
	a.activePlayer = player
	a.mu.Unlock()

	// Emit event to frontend
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "media:update", player)
	}
}

// LoadWindowPosition returns saved window position
func (a *App) LoadWindowPosition() (int, int) {
	return a.config.WindowX, a.config.WindowY
}

// SaveWindowPosition saves current window position
func (a *App) SaveWindowPosition() {
	if a.ctx == nil {
		return
	}
	x, y := runtime.WindowGetPosition(a.ctx)
	// Only save if position changed
	if x != a.config.WindowX || y != a.config.WindowY {
		a.config.WindowX = x
		a.config.WindowY = y
		a.config.Save()
	}
}

// startPositionSaver periodically saves window position
func (a *App) startPositionSaver() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.SaveWindowPosition()
		}
	}
}

// GetCurrentPlayer returns the current active player state
func (a *App) GetCurrentPlayer() *media.Player {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.activePlayer
}

// --- Media Control Methods ---

// MediaPlay sends play command to active player
func (a *App) MediaPlay() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "STATE", 0) // 0 = PLAYING
}

// MediaPause sends pause command to active player
func (a *App) MediaPause() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "STATE", 1) // 1 = PAUSED
}

// MediaTogglePlayPause toggles play/pause state
func (a *App) MediaTogglePlayPause() error {
	if a.activePlayer == nil {
		return nil
	}
	if a.activePlayer.State == media.StatePlaying {
		return a.MediaPause()
	}
	return a.MediaPlay()
}

// MediaNext sends next track command
func (a *App) MediaNext() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "SKIP_NEXT", nil)
}

// MediaPrevious sends previous track command
func (a *App) MediaPrevious() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "SKIP_PREVIOUS", nil)
}

// MediaToggleShuffle toggles shuffle mode
func (a *App) MediaToggleShuffle() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	newState := !a.activePlayer.Shuffle
	var val int
	if newState {
		val = 1
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "SHUFFLE", val)
}

// MediaToggleRepeat cycles through repeat modes
func (a *App) MediaToggleRepeat() error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	// Cycle: NONE(1) -> ALL(2) -> ONE(4) -> NONE(1)
	var nextMode int
	switch a.activePlayer.Repeat {
	case media.RepeatNone:
		nextMode = int(media.RepeatAll)
	case media.RepeatAll:
		nextMode = int(media.RepeatOne)
	default:
		nextMode = int(media.RepeatNone)
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "REPEAT", nextMode)
}

// MediaSetRating sets track rating (0=none, 1=dislike, 5=like)
func (a *App) MediaSetRating(rating int) error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "RATING", rating)
}

// MediaSeek seeks to position in seconds
func (a *App) MediaSeek(position int) error {
	if a.wnpServer == nil || a.activePlayer == nil {
		return nil
	}
	return a.wnpServer.SendCommand(a.activePlayer.ID, "POSITION", position)
}
