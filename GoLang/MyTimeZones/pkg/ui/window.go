package ui

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/MyTimeZones/pkg/config"
	"github.com/yourusername/MyTimeZones/pkg/logger"
	"github.com/yourusername/MyTimeZones/pkg/timezone"
)

type Window struct {
	app         fyne.App
	window      fyne.Window
	timeManager *timezone.Manager
	logger      *logger.Logger
	refreshRate time.Duration
	showSeconds bool
	ctx         context.Context
	cancel      context.CancelFunc
	statusBar   *widget.Label
	table       *widget.Table
	editWindow  *EditZonesWindow
	config      *config.AppConfig
}

func NewWindow(app fyne.App, config *config.AppConfig, timeManager *timezone.Manager, logger *logger.Logger, refreshRateSeconds int, showSeconds bool) *Window {
	if refreshRateSeconds < 1 {
		refreshRateSeconds = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Window{
		app:         app,
		config:      config,
		timeManager: timeManager,
		logger:      logger,
		refreshRate: time.Duration(refreshRateSeconds) * time.Second,
		showSeconds: showSeconds,
		ctx:         ctx,
		cancel:      cancel,
		statusBar:   widget.NewLabel(""),
	}
}

func (w *Window) Show() {
	w.window = w.app.NewWindow("MyTime")
	w.setupUI()
	w.startRefreshTimer()
	w.window.Show()
}

func (w *Window) setupUI() {
	w.table = w.createTimeTable()
	content := container.NewBorder(
		widget.NewLabel(""), // top
		w.statusBar,         // bottom
		nil, nil,            // left, right
		w.table, // center
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(800, 600))
	w.setupMenu()
}

func (w *Window) createTimeTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			timeInfo, _ := w.timeManager.GetTimeInfo()
			return len(timeInfo) + 1, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				headers := []string{"Name", "Description", "Date", "Time", "HoursDiff"}
				label.SetText(headers[i.Col])
				return
			}

			timeInfo, err := w.timeManager.GetTimeInfo()
			if err != nil {
				w.logger.Error("Failed to get time info: %v", err)
				return
			}

			if i.Row-1 < len(timeInfo) {
				info := timeInfo[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(info.Name)
				case 1:
					label.SetText(info.Description)
				case 2:
					label.SetText(info.Date)
				case 3:
					loc, _ := time.LoadLocation(info.Name)
					t := time.Now().In(loc)
					var timeFormat string
					if w.showSeconds {
						timeFormat = "15:04:05"
					} else {
						timeFormat = "15:04"
					}
					label.SetText(t.Format(timeFormat))
				case 4:
					label.SetText(info.Diff)
				}
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 120)
	table.SetColumnWidth(4, 100)

	return table
}

func (w *Window) startRefreshTimer() {
	go func() {
		ticker := time.NewTicker(w.refreshRate)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				w.refresh()
			}
		}
	}()
}

func (w *Window) refresh() {
	w.table.Refresh()
	w.statusBar.SetText("Last updated: " + time.Now().Format("15:04:05"))
}

func (w *Window) setupMenu() {
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Add Timezone", func() {
				addWindow := NewAddZonesWindow(w.app, w.config, w.timeManager)
				addWindow.Show()
			}),
			fyne.NewMenuItem("Edit Zones", w.showEditZonesWindow),
			fyne.NewMenuItem("Close", w.close),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("About", w.showAbout),
		),
	)
	w.window.SetMainMenu(mainMenu)
}

func (w *Window) showTimeWindow() {
	newWindow := w.app.NewWindow("Show Time")
	newWindow.SetContent(w.createTimeTable())
	newWindow.Resize(fyne.NewSize(800, 600))
	newWindow.Show()
}

func (w *Window) showEditZonesWindow() {
	if w.editWindow == nil {
		w.editWindow = NewEditZonesWindow(w.app, w.config, w.timeManager)
	}
	w.editWindow.Show()
}

func (w *Window) showAbout() {
	dialog.ShowInformation("About", "MyTime - Time Zone Manager", w.window)
}

func (w *Window) close() {
	w.cancel()
	w.app.Quit()
}
