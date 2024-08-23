package tests

import (
	"os"
	"testing"
	"wakey/internal/config"
)

func TestReadConfig(t *testing.T) {
	// Setup: Create a temporary config file
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write sample config data to the temp file
	sampleConfig := `{"Devices": [{"DeviceName": "Device1", "Description": "This is a test device.", "MACAddress": "00:00:00:00:00:00", "IPADdress": "1.1.1.1"}]}`
	if _, err := tempFile.Write([]byte(sampleConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the ConfigPath to point to the temp file
	config.ConfigPath = tempFile.Name()

	// Execute: Call ReadConfig
	cfg := config.ReadConfig()

	// Verify: Check if the config was read correctly
	if len(cfg.Devices) != 1 || cfg.Devices[0].DeviceName != "Device1" {
		t.Errorf("Expected Device1, got %v", cfg.Devices)
	}
}
