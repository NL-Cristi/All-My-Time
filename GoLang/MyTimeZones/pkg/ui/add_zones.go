package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/MyTimeZones/pkg/config"
	"github.com/yourusername/MyTimeZones/pkg/timezone"
)

type AddZonesWindow struct {
	app           fyne.App
	window        fyne.Window
	config        *config.AppConfig
	timeManager   *timezone.Manager
	searchEntry   *widget.Entry
	description   *widget.Entry
	zonesList     *widget.List
	filteredZones []string
	selectedIndex int
	allZones      []string // Add this field to store all timezones
}

func NewAddZonesWindow(app fyne.App, config *config.AppConfig, timeManager *timezone.Manager) *AddZonesWindow {
	allZones := timezone.GetTimeZones()
	return &AddZonesWindow{
		app:           app,
		config:        config,
		timeManager:   timeManager,
		filteredZones: make([]string, len(allZones)),
		allZones:      allZones,
		selectedIndex: -1,
	}
}

func (a *AddZonesWindow) Show() {
	if a.window == nil {
		a.window = a.app.NewWindow("Add Timezone")
		a.createUI()
		a.window.Resize(fyne.NewSize(500, 600))
	}
	a.window.Show()
}

func (a *AddZonesWindow) createUI() {
	// Search field
	a.searchEntry = widget.NewEntry()
	a.searchEntry.SetPlaceHolder("Search timezones...")
	a.searchEntry.OnChanged = a.filterZones

	// Description field
	a.description = widget.NewEntry()
	a.description.SetPlaceHolder("Enter description...")

	// Initialize filtered zones with all timezones
	copy(a.filteredZones, a.allZones)

	// Timezone list
	a.zonesList = widget.NewList(
		func() int {
			return len(a.filteredZones)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(a.filteredZones[id])
		},
	)

	// Add selection handler
	a.zonesList.OnSelected = func(id widget.ListItemID) {
		a.selectedIndex = int(id)
	}

	// The scrollable list (will take all remaining space)
	listScroll := container.NewVScroll(a.zonesList)

	// Add button
	addButton := widget.NewButton("Add Selected Timezone", func() {
		a.addSelectedTimezone()
	})

	// Top part: everything above the list
	topContent := container.NewVBox(
		widget.NewLabel("Search Timezones"),
		a.searchEntry,
		widget.NewLabel("Description"),
		a.description,
		addButton,
	)

	// Use Border layout: top is your form, center is the list
	content := container.NewBorder(
		topContent, // top
		nil,        // bottom
		nil, nil,   // left, right
		listScroll, // center
	)
	a.window.SetContent(content)
}

func (a *AddZonesWindow) filterZones(searchText string) {
	searchText = strings.ToLower(searchText)
	a.filteredZones = nil

	if searchText == "" {
		// If search is empty, show all timezones
		a.filteredZones = make([]string, len(a.allZones))
		copy(a.filteredZones, a.allZones)
	} else {
		// Filter timezones based on search text
		for _, tz := range a.allZones {
			if strings.Contains(strings.ToLower(tz), searchText) {
				a.filteredZones = append(a.filteredZones, tz)
			}
		}
	}

	a.zonesList.Refresh()
}

func (a *AddZonesWindow) addSelectedTimezone() {
	if a.selectedIndex == -1 {
		dialog.ShowError(fmt.Errorf("please select a timezone"), a.window)
		return
	}

	if a.description.Text == "" {
		dialog.ShowError(fmt.Errorf("please enter a description"), a.window)
		return
	}

	a.config.TimeZones.Others = append(a.config.TimeZones.Others, timezone.TimeZoneEntry{
		Zone:        a.filteredZones[a.selectedIndex],
		Description: a.description.Text,
	})

	if err := a.config.Save("config.json"); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save timezone: %v", err), a.window)
		return
	}
	if err := a.timeManager.UpdateConfig(a.config.TimeZones); err != nil {
		dialog.ShowError(fmt.Errorf("failed to update manager: %v", err), a.window)
		return
	}

	dialog.ShowInformation("Success", "Timezone added successfully", a.window)
	a.window.Close()
}
