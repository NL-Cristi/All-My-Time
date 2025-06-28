package timezone

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantErr    bool
	}{
		{
			name:       "Valid config file",
			configFile: "testdata/valid_config.json",
			wantErr:    false,
		},
		{
			name:       "Invalid config file",
			configFile: "testdata/invalid_config.json",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManager(tt.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetTimeInfo(t *testing.T) {
	manager, err := NewManager("testdata/valid_config.json")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	timeInfo, err := manager.GetTimeInfo()
	if err != nil {
		t.Fatalf("GetTimeInfo() error = %v", err)
	}

	if len(timeInfo) == 0 {
		t.Error("GetTimeInfo() returned empty slice")
	}

	// Verify local timezone is first
	if timeInfo[0].Diff != "00:00" {
		t.Errorf("Local timezone diff should be 00:00, got %s", timeInfo[0].Diff)
	}
}
