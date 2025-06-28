package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/yourusername/MyTimeZones/pkg/config"
	"github.com/yourusername/MyTimeZones/pkg/timezone"
)

type EditZonesWindow struct {
	app                fyne.App
	window             fyne.Window
	timeManager        *timezone.Manager
	localZone          *widget.Entry
	localDesc          *widget.Entry
	otherZones         *widget.List
	config             *config.AppConfig
	selectedOtherIndex int // -1 means editing local zone
}

func NewEditZonesWindow(app fyne.App, config *config.AppConfig, timeManager *timezone.Manager) *EditZonesWindow {
	return &EditZonesWindow{
		app:         app,
		config:      config,
		timeManager: timeManager,
	}
}

func (e *EditZonesWindow) Show() {
	// Always reload the config to get the latest data
	e.config.TimeZones = e.timeManager.GetConfig()
	e.selectedOtherIndex = -1

	if e.window == nil {
		e.window = e.app.NewWindow("Edit Time Zones")
		e.createUI()
		e.window.Resize(fyne.NewSize(600, 500))
		// Ensure the window can be recreated after closing
		e.window.SetOnClosed(func() {
			e.window = nil
		})
	} else {
		// If window already exists, just refresh the list and fields
		e.otherZones.Refresh()
		e.localZone.SetText(e.config.TimeZones.Local.Zone)
		e.localDesc.SetText(e.config.TimeZones.Local.Description)
	}
	e.window.Show()
}

func (e *EditZonesWindow) createUI() {
	// Local timezone section
	e.localZone = widget.NewEntry()
	e.localDesc = widget.NewEntry()
	e.localZone.SetText(e.config.TimeZones.Local.Zone)
	e.localDesc.SetText(e.config.TimeZones.Local.Description)

	localForm := widget.NewForm(
		widget.NewFormItem("Local Zone", e.localZone),
		widget.NewFormItem("Description", e.localDesc),
	)

	// Other timezones section as a list
	e.selectedOtherIndex = -1
	e.otherZones = widget.NewList(
		func() int {
			return len(e.config.TimeZones.Others)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewButton("Remove", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			button := box.Objects[1].(*widget.Button)

			tz := e.config.TimeZones.Others[id]
			label.SetText(fmt.Sprintf("%s - %s", tz.Zone, tz.Description))
			button.SetText("Remove")
			button.OnTapped = func() {
				e.removeZone(id)
				e.selectedOtherIndex = -1
				e.localZone.SetText(e.config.TimeZones.Local.Zone)
				e.localDesc.SetText(e.config.TimeZones.Local.Description)
			}
		},
	)
	e.otherZones.OnSelected = func(id widget.ListItemID) {
		e.selectedOtherIndex = id
		tz := e.config.TimeZones.Others[id]
		e.localZone.SetText(tz.Zone)
		e.localDesc.SetText(tz.Description)
	}

	// The scrollable list (will take all remaining space)
	otherZonesScroll := container.NewVScroll(e.otherZones)

	// Add button to open the add timezone window
	addButton := widget.NewButton("Add New Timezone", func() {
		addWindow := NewAddZonesWindow(e.app, e.config, e.timeManager)
		addWindow.Show()
	})

	// Save button
	saveButton := widget.NewButton("Save Changes", func() {
		e.saveChanges()
	})

	// Create a compact button container
	buttonContainer := container.NewHBox(
		addButton,
		layout.NewSpacer(),
		saveButton,
	)

	// Top part: everything above the list
	topContent := container.NewVBox(
		buttonContainer,
		widget.NewSeparator(),
		widget.NewLabel("Local Timezone"),
		localForm,
		widget.NewSeparator(),
		widget.NewLabel("Other Timezones"),
	)

	// Use Border layout: top is your form, center is the list
	content := container.NewBorder(
		topContent, // top
		nil,        // bottom
		nil, nil,   // left, right
		otherZonesScroll, // center
	)
	e.window.SetContent(content)
}

func (e *EditZonesWindow) removeZone(index int) {
	e.config.TimeZones.Others = append(e.config.TimeZones.Others[:index], e.config.TimeZones.Others[index+1:]...)
	_ = e.config.Save("config.json")
	_ = e.timeManager.UpdateConfig(e.config.TimeZones)
	e.otherZones.Refresh()
}

func (e *EditZonesWindow) saveChanges() {
	// Validate the timezone
	if _, err := time.LoadLocation(e.localZone.Text); err != nil {
		dialog.ShowError(fmt.Errorf("invalid timezone: %v", err), e.window)
		return
	}

	if e.selectedOtherIndex >= 0 && e.selectedOtherIndex < len(e.config.TimeZones.Others) {
		// Update selected other zone
		e.config.TimeZones.Others[e.selectedOtherIndex].Zone = e.localZone.Text
		e.config.TimeZones.Others[e.selectedOtherIndex].Description = e.localDesc.Text
	} else {
		// Update local zone
		e.config.TimeZones.Local.Zone = e.localZone.Text
		e.config.TimeZones.Local.Description = e.localDesc.Text
	}

	if err := e.config.Save("config.json"); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save configuration: %v", err), e.window)
		return
	}
	if err := e.timeManager.UpdateConfig(e.config.TimeZones); err != nil {
		dialog.ShowError(fmt.Errorf("failed to update manager: %v", err), e.window)
		return
	}

	dlg := dialog.NewInformation("Success", "Timezone configuration saved", e.window)
	dlg.SetOnClosed(func() {
		e.window.Close()
	})
	dlg.Show()
}
