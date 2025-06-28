package main

import (
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

type TimeZoneConfig struct {
	Local struct {
		Zone        string `json:"zone"`
		Description string `json:"description"`
	} `json:"local"`
	Others []struct {
		Zone        string `json:"zone"`
		Description string `json:"description"`
	} `json:"others"`
}

type TimeZoneManager struct {
	config     TimeZoneConfig
	configFile string
}

func NewTimeZoneManager(configFile string) (*TimeZoneManager, error) {
	tzm := &TimeZoneManager{
		configFile: configFile,
	}

	if err := tzm.loadConfig(); err != nil {
		return nil, err
	}

	return tzm, nil
}

func (tzm *TimeZoneManager) loadConfig() error {
	file, err := os.ReadFile(tzm.configFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &tzm.config)
}

func (tzm *TimeZoneManager) saveConfig() error {
	data, err := json.MarshalIndent(tzm.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tzm.configFile, data, 0644)
}

func (tzm *TimeZoneManager) GetTimeInfo() []TimeInfo {
	localLoc, _ := time.LoadLocation(tzm.config.Local.Zone)
	localTime := time.Now().In(localLoc)

	timeInfo := make([]TimeInfo, 0)

	// Add local timezone info
	timeInfo = append(timeInfo, TimeInfo{
		Name:        tzm.config.Local.Zone,
		Description: tzm.config.Local.Description,
		Date:        localTime.Format("2006-01-02"),
		Time:        localTime.Format("15:04"),
		Diff:        "00:00",
	})

	// Add other timezone info
	for _, tz := range tzm.config.Others {
		loc, _ := time.LoadLocation(tz.Zone)
		currentTime := time.Now().In(loc)

		// Calculate time zone offset difference
		_, localOffset := localTime.Zone()
		_, otherOffset := currentTime.Zone()
		offsetDiff := otherOffset - localOffset
		sign := "+"
		if offsetDiff < 0 {
			sign = "-"
			offsetDiff = -offsetDiff
		}
		hours := offsetDiff / 3600
		minutes := (offsetDiff % 3600) / 60

		timeInfo = append(timeInfo, TimeInfo{
			Name:        tz.Zone,
			Description: tz.Description,
			Date:        currentTime.Format("2006-01-02"),
			Time:        currentTime.Format("15:04"),
			Diff:        fmt.Sprintf("%s%02d:%02d", sign, hours, minutes),
		})
	}

	return timeInfo
}

func formatTimeDiff(hours, minutes int) string {
	sign := "+"
	if hours < 0 {
		sign = "-"
		hours = -hours
	}
	return sign + formatNumber(hours) + ":" + formatNumber(minutes)
}

func formatNumber(n int) string {
	if n < 10 {
		return "0" + string(rune(n+'0'))
	}
	return string(rune(n/10+'0')) + string(rune(n%10+'0'))
}
