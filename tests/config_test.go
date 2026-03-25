package tests

import (
	"os"
	"testing"
	"wakey/internal/config"
)

func TestReadConfig(t *testing.T) {
	// Setup: use a temp SQLite database
	tmpFile, err := os.CreateTemp("", "wakey_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp db file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	config.DBPath = tmpFile.Name()
	config.ResetDB()
	if err := config.CreateConfig(); err != nil {
		// CreateConfig returns a non-nil error on success (used as status message)
		_ = err
	}

	// Write a device directly via WriteConfig
	sample := config.Config{
		Devices: []config.Device{
			{
				ID:          "test-id-1",
				DeviceName:  "Device1",
				Description: "This is a test device.",
				MacAddress:  "00:00:00:00:00:00",
				IPAddress:   "1.1.1.1",
				State:       "Offline",
			},
		},
		Groups: []config.Group{},
	}
	config.WriteConfig(sample)

	// Execute: Call ReadConfig
	cfg := config.ReadConfig()

	// Verify
	if len(cfg.Devices) != 1 || cfg.Devices[0].DeviceName != "Device1" {
		t.Errorf("Expected Device1, got %v", cfg.Devices)
	}
}
