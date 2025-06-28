package main

import (
	"os"

	"fyne.io/fyne/v2/app"

	"github.com/yourusername/GoMyTime/pkg/config"
	"github.com/yourusername/GoMyTime/pkg/logger"
	"github.com/yourusername/GoMyTime/pkg/timezone"
	"github.com/yourusername/GoMyTime/pkg/ui"
)

func main() {
	log := logger.NewLogger("info")

	cfg, err := config.LoadOrCreateConfig("config.json")
	if err != nil {
		log.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	refreshRate := cfg.RefreshRateSeconds
	if !cfg.ShowSeconds && refreshRate < 60 {
		refreshRate = 60
	}

	timeManager, err := timezone.NewManager("config.json")
	if err != nil {
		log.Error("Failed to initialize timezone manager: %v", err)
		os.Exit(1)
	}

	// Ensure manager is in sync with config
	_ = timeManager.UpdateConfig(cfg.TimeZones)

	myApp := app.New()
	window := ui.NewWindow(myApp, cfg, timeManager, log, refreshRate, cfg.ShowSeconds)
	window.Show()
	myApp.Run()
}
