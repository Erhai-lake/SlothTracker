package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "SlothTracker",
		Width:             1024,
		Height:            768,
		HideWindowOnClose: true,
		DisableResize:     false,
		Frameless:         false,
		StartHidden:       false,
		Bind: []any{
			app,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
