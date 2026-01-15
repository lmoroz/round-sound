package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"round-sound/app"
)

//go:embed all:frontend/dist
var assets embed.FS

// getWebViewDataPath returns the path for WebView2 user data
func getWebViewDataPath() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	// Use the same folder as config (without .exe suffix)
	dataDir := filepath.Join(appData, "round-sound", "webview")
	os.MkdirAll(dataDir, 0755)
	return dataDir
}

func main() {
	// Create application instance
	application := app.NewApp()

	// Load saved window position
	x, y := application.LoadWindowPosition()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Round Sound",
		Width:             600,
		Height:            600,
		MinWidth:          600,
		MinHeight:         600,
		MaxWidth:          600,
		MaxHeight:         600,
		DisableResize:     true,
		Frameless:         true,
		StartHidden:       false,
		HideWindowOnClose: true,
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: true,
			WebviewUserDataPath:               getWebViewDataPath(),
			ZoomFactor:                        1.0,
		},
		OnStartup:  application.Startup,
		OnShutdown: application.Shutdown,
		OnDomReady: application.DomReady,
		Bind: []interface{}{
			application,
		},
	})

	// Set initial window position if saved
	if x != 0 || y != 0 {
		// Position will be set in DomReady
	}

	if err != nil {
		log.Fatal(err)
	}
}
