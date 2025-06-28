package timezone

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TimeInfo struct {
	Name        string
	Description string
	Date        string
	Time        string
	Diff        string
}

// TimeZoneConfig represents the configuration structure for timezones
type TimeZoneConfig struct {
	Local  TimeZoneEntry   `json:"local"`
	Others []TimeZoneEntry `json:"others"`
}

// TimeZoneEntry represents a single timezone entry
type TimeZoneEntry struct {
	Zone        string `json:"zone"`
	Description string `json:"description"`
}

var DefaultTimeZoneConfig = TimeZoneConfig{
	Local: TimeZoneEntry{
		Zone:        "Europe/Lisbon",
		Description: "Local Time",
	},
	Others: []TimeZoneEntry{
		{Zone: "Europe/Bucharest", Description: "Bucharest Time"},
		{Zone: "America/New_York", Description: "New York Time"},
		{Zone: "Europe/London", Description: "London Time"},
		{Zone: "Europe/Paris", Description: "Paris Time"},
		{Zone: "Europe/Berlin", Description: "Berlin Time"},
		{Zone: "Asia/Tokyo", Description: "Tokyo Time"},
		{Zone: "Australia/Sydney", Description: "Sydney Time"},
		{Zone: "America/Los_Angeles", Description: "Los Angeles Time"},
		{Zone: "Asia/Kolkata", Description: "Kolkata Time"},
	},
}

type Manager struct {
	config     TimeZoneConfig
	configFile string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewManager(configFile string) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	tzm := &Manager{
		configFile: configFile,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Always load the time zone section from the file
	_ = tzm.loadConfig()

	return tzm, nil
}

func (m *Manager) validateConfig() error {
	// Validate local timezone
	if _, err := time.LoadLocation(m.config.Local.Zone); err != nil {
		return fmt.Errorf("invalid local timezone %s: %w", m.config.Local.Zone, err)
	}

	// Validate other timezones
	for _, tz := range m.config.Others {
		if _, err := time.LoadLocation(tz.Zone); err != nil {
			return fmt.Errorf("invalid timezone %s: %w", tz.Zone, err)
		}
	}

	return nil
}

func (m *Manager) loadConfig() error {
	file, err := os.ReadFile(m.configFile)
	if err != nil {
		return err
	}
	// Only unmarshal the timeZones section
	var fileData map[string]interface{}
	if err := json.Unmarshal(file, &fileData); err != nil {
		return err
	}
	if tzSection, ok := fileData["timeZones"]; ok {
		tzBytes, _ := json.Marshal(tzSection)
		return json.Unmarshal(tzBytes, &m.config)
	}
	return nil
}

func (m *Manager) createDefaultConfig() error {
	m.config.Local = TimeZoneEntry{
		Zone:        "UTC",
		Description: "Coordinated Universal Time",
	}
	return nil
}

func (m *Manager) GetTimeInfo() ([]TimeInfo, error) {
	localLoc, err := time.LoadLocation(m.config.Local.Zone)
	if err != nil {
		return nil, fmt.Errorf("failed to load local timezone: %w", err)
	}

	localTime := time.Now().In(localLoc)
	timeInfo := make([]TimeInfo, 0, len(m.config.Others)+1)

	// Add local timezone info
	timeInfo = append(timeInfo, TimeInfo{
		Name:        m.config.Local.Zone,
		Description: m.config.Local.Description,
		Date:        localTime.Format("2006-01-02"),
		Time:        localTime.Format("15:04:05"),
		Diff:        "00:00",
	})

	// Add other timezone info
	for _, tz := range m.config.Others {
		loc, err := time.LoadLocation(tz.Zone)
		if err != nil {
			return nil, fmt.Errorf("failed to load timezone %s: %w", tz.Zone, err)
		}

		currentTime := time.Now().In(loc)
		_, localOffset := localTime.Zone()
		_, otherOffset := currentTime.Zone()
		offsetDiff := otherOffset - localOffset

		timeInfo = append(timeInfo, TimeInfo{
			Name:        tz.Zone,
			Description: tz.Description,
			Date:        currentTime.Format("2006-01-02"),
			Time:        currentTime.Format("15:04:05"),
			Diff:        formatOffset(offsetDiff),
		})
	}

	return timeInfo, nil
}

func formatOffset(offsetSeconds int) string {
	sign := "+"
	if offsetSeconds < 0 {
		sign = "-"
		offsetSeconds = -offsetSeconds
	}

	hours := offsetSeconds / 3600
	minutes := (offsetSeconds % 3600) / 60

	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

func (m *Manager) Close() {
	if m.cancel != nil {
		m.cancel()
	}
}

// GetConfig returns a copy of the current configuration
func (m *Manager) GetConfig() TimeZoneConfig {
	return m.config
}

// UpdateConfig updates the timezone configuration and validates it
func (m *Manager) UpdateConfig(config TimeZoneConfig) error {
	m.config = config
	if err := m.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	return nil
}

func NewManagerFromConfig(cfg TimeZoneConfig) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	tzm := &Manager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
	if err := tzm.validateConfig(); err != nil {
		return nil, err
	}
	return tzm, nil
}
