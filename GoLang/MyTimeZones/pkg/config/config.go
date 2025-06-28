package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/MyTimeZones/pkg/timezone"
)

type AppConfig struct {
	WindowWidth        int                     `json:"windowWidth"`
	WindowHeight       int                     `json:"windowHeight"`
	RefreshRateSeconds int                     `json:"refreshRateSeconds"`
	ShowSeconds        bool                    `json:"showSeconds"`
	TimeZones          timezone.TimeZoneConfig `json:"timeZones"`
}

var DefaultAppConfig = AppConfig{
	WindowWidth:        800,
	WindowHeight:       600,
	RefreshRateSeconds: 1,
	ShowSeconds:        true,
	TimeZones:          timezone.DefaultTimeZoneConfig,
}

func LoadOrCreateConfig(path string) (*AppConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File does not exist, create it with default config
		data, _ := json.MarshalIndent(DefaultAppConfig, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.RefreshRateSeconds < 1 {
		cfg.RefreshRateSeconds = 1
	}

	return &cfg, nil
}

func (c *AppConfig) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
