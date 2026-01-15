package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"round-sound/media"
)

// App struct holds the application state
type App struct {
	ctx            context.Context
	mu             sync.RWMutex
	config         *Config
	wnpServer      *media.WebNowPlayingServer
	activePlayer   *media.Player
	windowManager  *WindowManager
	audioCapture   *media.AudioLevelCapture
	trayManager    *TrayManager
	autorunManager *AutorunManager
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

	// Initialize and setup system tray
	a.trayManager = NewTrayManager(ctx)
	a.trayManager.Setup()

	// Initialize autorun manager
	var err error
	a.autorunManager, err = NewAutorunManager()
	if err != nil {
		log.Printf("Failed to initialize autorun manager: %v", err)
	}

	// Start WebNowPlaying server on configured port
	a.startWNPServer(a.config.WNPPort)

	// Start desktop-level window manager (HWND_BOTTOM)
	go a.windowManager.StartDesktopLevelWatcher()

	// Start audio level capture
	a.audioCapture = media.NewAudioLevelCapture(a.onAudioLevels)
	if err = a.audioCapture.Start(); err != nil {
		log.Printf("Failed to start audio capture: %v", err)
	}

	// Listen for audio configuration changes from frontend
	runtime.EventsOn(ctx, "audio:config", a.onAudioConfigUpdate)

	log.Println("Round Sound started")
}

// DomReady is called after the front-end dom has been loaded
func (a *App) DomReady(ctx context.Context) {
	// Set saved window position
	if a.config.WindowX != 0 || a.config.WindowY != 0 {
		runtime.WindowSetPosition(ctx, a.config.WindowX, a.config.WindowY)
		log.Printf("Restored window position: %d, %d", a.config.WindowX, a.config.WindowY)
	}

	// Set window tode desktop level
	a.windowManager.SetDesktopLevel()

	// Hide from taskbar (will only show in system tray)
	if a.windowManager.Hwnd != 0 {
		setToolWindow(a.windowManager.Hwnd)
		log.Println("Window hidden from taskbar")
	}

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

	// Stop audio capture
	if a.audioCapture != nil {
		a.audioCapture.Stop()
	}

	// Remove tray icon
	if a.trayManager != nil {
		a.trayManager.Remove()
	}

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

	log.Printf("[App] Player updated: ID=%d, Title=%s, State=%d", player.ID, player.Title, player.State)

	// Emit event to frontend
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "media:update", player)
	}
}

// onAudioLevels is called when audio levels are captured
func (a *App) onAudioLevels(levels []float32) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "audio:levels", levels)
	}
}

// onAudioConfigUpdate is called when frontend sends audio configuration changes
func (a *App) onAudioConfigUpdate(data ...interface{}) {
	if len(data) == 0 {
		return
	}

	config, ok := data[0].(map[string]interface{})
	if !ok {
		log.Printf("[App] Invalid audio config format")
		return
	}

	fftSize, _ := config["fftSize"].(float64)
	freqMin, _ := config["freqMin"].(float64)
	freqMax, _ := config["freqMax"].(float64)

	if a.audioCapture != nil {
		a.audioCapture.UpdateConfig(int(fftSize), freqMin, freqMax)
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
	log.Println("[App] MediaPlay called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaPlay: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaPlay: no active player")
		return nil
	}

	log.Printf("[App] Sending STATE command (PLAYING) to player %d", a.activePlayer.ID)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "STATE", 1) // 1 = PLAYING
	if err != nil {
		log.Printf("[App] MediaPlay error: %v", err)
	}
	return err
}

// MediaPause sends pause command to active player
func (a *App) MediaPause() error {
	log.Println("[App] MediaPause called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaPause: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaPause: no active player")
		return nil
	}

	log.Printf("[App] Sending STATE command (PAUSED) to player %d", a.activePlayer.ID)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "STATE", 2) // 2 = PAUSED
	if err != nil {
		log.Printf("[App] MediaPause error: %v", err)
	}
	return err
}

// MediaTogglePlayPause toggles play/pause state
func (a *App) MediaTogglePlayPause() error {
	log.Println("[App] MediaTogglePlayPause called")

	a.mu.RLock()
	currentState := media.StateStopped
	if a.activePlayer != nil {
		currentState = a.activePlayer.State
	}
	a.mu.RUnlock()

	if a.activePlayer == nil {
		log.Println("[App] MediaTogglePlayPause: no active player")
		return nil
	}

	if currentState == media.StatePlaying {
		return a.MediaPause()
	}
	return a.MediaPlay()
}

// MediaNext sends next track command
func (a *App) MediaNext() error {
	log.Println("[App] MediaNext called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaNext: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaNext: no active player")
		return nil
	}

	log.Printf("[App] Sending SKIP_NEXT command to player %d", a.activePlayer.ID)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "SKIP_NEXT", nil)
	if err != nil {
		log.Printf("[App] MediaNext error: %v", err)
	}
	return err
}

// MediaPrevious sends previous track command
func (a *App) MediaPrevious() error {
	log.Println("[App] MediaPrevious called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaPrevious: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaPrevious: no active player")
		return nil
	}

	log.Printf("[App] Sending SKIP_PREVIOUS command to player %d", a.activePlayer.ID)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "SKIP_PREVIOUS", nil)
	if err != nil {
		log.Printf("[App] MediaPrevious error: %v", err)
	}
	return err
}

// MediaToggleShuffle toggles shuffle mode
func (a *App) MediaToggleShuffle() error {
	log.Println("[App] MediaToggleShuffle called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaToggleShuffle: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaToggleShuffle: no active player")
		return nil
	}

	newState := !a.activePlayer.Shuffle
	var val int
	if newState {
		val = 1
	}
	log.Printf("[App] Sending SHUFFLE command to player %d (newState=%v)", a.activePlayer.ID, newState)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "SHUFFLE", val)
	if err != nil {
		log.Printf("[App] MediaToggleShuffle error: %v", err)
	}
	return err
}

// MediaToggleRepeat cycles through repeat modes
func (a *App) MediaToggleRepeat() error {
	log.Println("[App] MediaToggleRepeat called")

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaToggleRepeat: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaToggleRepeat: no active player")
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
	log.Printf("[App] Sending REPEAT command to player %d (nextMode=%d)", a.activePlayer.ID, nextMode)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "REPEAT", nextMode)
	if err != nil {
		log.Printf("[App] MediaToggleRepeat error: %v", err)
	}
	return err
}

// MediaSetRating sets track rating (0=none, 1=dislike, 5=like)
func (a *App) MediaSetRating(rating int) error {
	log.Printf("[App] MediaSetRating called with rating=%d", rating)

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaSetRating: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaSetRating: no active player")
		return nil
	}

	log.Printf("[App] Sending RATING command to player %d (rating=%d)", a.activePlayer.ID, rating)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "RATING", rating)
	if err != nil {
		log.Printf("[App] MediaSetRating error: %v", err)
	}
	return err
}

// MediaSeek seeks to position in seconds
func (a *App) MediaSeek(position int) error {
	log.Printf("[App] MediaSeek called with position=%d", position)

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaSeek: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaSeek: no active player")
		return nil
	}

	log.Printf("[App] Sending POSITION command to player %d (position=%d)", a.activePlayer.ID, position)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "POSITION", position)
	if err != nil {
		log.Printf("[App] MediaSeek error: %v", err)
	}
	return err
}

// MediaSetVolume sets the volume (0-100)
func (a *App) MediaSetVolume(volume int) error {
	// Clamp volume between 0 and 100
	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}

	log.Printf("[App] MediaSetVolume called with volume=%d", volume)

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.wnpServer == nil {
		log.Println("[App] MediaSetVolume: wnpServer is nil")
		return nil
	}
	if a.activePlayer == nil {
		log.Println("[App] MediaSetVolume: no active player")
		return nil
	}

	if !a.activePlayer.CanSetVolume {
		log.Println("[App] MediaSetVolume: player does not support setting volume")
		return nil
	}

	log.Printf("[App] Sending VOLUME command to player %d (volume=%d)", a.activePlayer.ID, volume)
	err := a.wnpServer.SendCommand(a.activePlayer.ID, "VOLUME", volume)
	if err != nil {
		log.Printf("[App] MediaSetVolume error: %v", err)
	}
	return err
}

// --- Autorun Methods ---

// IsAutorunEnabled checks if autorun is enabled
func (a *App) IsAutorunEnabled() bool {
	if a.autorunManager == nil {
		return false
	}
	enabled, err := a.autorunManager.IsEnabled()
	if err != nil {
		log.Printf("[App] IsAutorunEnabled error: %v", err)
		return false
	}
	return enabled
}

// SetAutorun enables or disables autorun
func (a *App) SetAutorun(enabled bool) error {
	if a.autorunManager == nil {
		return nil
	}

	if enabled {
		return a.autorunManager.Enable()
	}
	return a.autorunManager.Disable()
}

// ShowWindow shows the window from tray
func (a *App) ShowWindow() {
	if a.trayManager != nil {
		a.trayManager.ShowWindow()
	}
}

// Quit gracefully exits the application
func (a *App) Quit() {
	log.Println("[App] Quit requested from frontend")
	if a.trayManager != nil {
		a.trayManager.Remove()
	}
	runtime.Quit(a.ctx)
}

// --- WebNowPlaying Port Management ---

// startWNPServer starts the WNP server on the specified port
func (a *App) startWNPServer(port int) {
	var err error
	a.wnpServer, err = media.NewWebNowPlayingServer(port, a.onPlayerUpdate)
	if err != nil {
		log.Printf("Failed to start WebNowPlaying server on port %d: %v", port, err)

		// Emit event to open settings with port configuration hint
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "wnp:port_busy", map[string]interface{}{
				"port":    port,
				"message": "Порт занят другим приложением (возможно, Rainmeter WebNowPlaying).",
			})
		}
	} else {
		log.Printf("WebNowPlaying server started on port %d", port)
	}
}

// GetWNPPort returns the current WNP port from config
func (a *App) GetWNPPort() int {
	if a.config != nil {
		return a.config.WNPPort
	}
	return DefaultWNPPort
}

// IsWNPConnected returns true if WNP server is running
func (a *App) IsWNPConnected() bool {
	return a.wnpServer != nil
}

// ChangeWNPPort changes the WNP port and restarts the server
func (a *App) ChangeWNPPort(port int) error {
	if port < 1024 || port > 65535 {
		return fmt.Errorf("invalid port: must be between 1024 and 65535")
	}

	log.Printf("Changing WNP port from %d to %d", a.config.WNPPort, port)

	// Stop existing server
	if a.wnpServer != nil {
		a.wnpServer.Stop()
		a.wnpServer = nil
	}

	// Update config
	a.config.WNPPort = port
	a.config.Save()

	// Start new server
	var err error
	a.wnpServer, err = media.NewWebNowPlayingServer(port, a.onPlayerUpdate)
	if err != nil {
		log.Printf("Failed to start WebNowPlaying server on port %d: %v", port, err)

		// Emit error event
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "wnp:port_error", map[string]interface{}{
				"port":    port,
				"message": "Не удалось запустить сервер на этом порту. Возможно, порт занят.",
			})
		}
		return err
	}

	log.Printf("WebNowPlaying server restarted on port %d", port)

	// Emit success event
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "wnp:port_changed", port)
	}

	return nil
}
