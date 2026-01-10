package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"round-sound/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create application instance
	application := app.NewApp()

	// Load saved window position
	x, y := application.LoadWindowPosition()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Round Sound",
		Width:             400,
		Height:            400,
		MinWidth:          400,
		MinHeight:         400,
		MaxWidth:          400,
		MaxHeight:         400,
		DisableResize:     true,
		Frameless:         true,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableWindowIcon:                 true,
			DisableFramelessWindowDecorations: true,
			WebviewUserDataPath:               "",
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
