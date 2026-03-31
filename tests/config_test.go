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

func TestWriteConfigMultipleDevices(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{
			{ID: "id-1", DeviceName: "Alpha", Description: "First", MacAddress: "AA:AA:AA:AA:AA:AA", IPAddress: "192.168.1.1", State: "Offline"},
			{ID: "id-2", DeviceName: "Beta", Description: "Second", MacAddress: "BB:BB:BB:BB:BB:BB", IPAddress: "192.168.1.2", State: "Offline"},
			{ID: "id-3", DeviceName: "Gamma", Description: "Third", MacAddress: "CC:CC:CC:CC:CC:CC", IPAddress: "192.168.1.3", State: "Offline"},
		},
		Groups: []config.Group{},
	}
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Devices) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(got.Devices))
	}
	names := []string{"Alpha", "Beta", "Gamma"}
	for i, want := range names {
		if got.Devices[i].DeviceName != want {
			t.Errorf("device[%d]: expected %q, got %q", i, want, got.Devices[i].DeviceName)
		}
	}
}

func TestGroupInsertionOrder(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{},
		Groups: []config.Group{
			{ID: "g-1", GroupName: "First", Devices: []string{}},
			{ID: "g-2", GroupName: "Second", Devices: []string{}},
			{ID: "g-3", GroupName: "Third", Devices: []string{}},
		},
	}
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(got.Groups))
	}
	names := []string{"First", "Second", "Third"}
	for i, want := range names {
		if got.Groups[i].GroupName != want {
			t.Errorf("group[%d]: expected %q, got %q", i, want, got.Groups[i].GroupName)
		}
	}
}

func TestGroupDeviceInsertionOrder(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{
			{ID: "d-1", DeviceName: "First", Description: "desc", MacAddress: "AA:AA:AA:AA:AA:AA", IPAddress: "192.168.1.1", State: "Offline"},
			{ID: "d-2", DeviceName: "Second", Description: "desc", MacAddress: "BB:BB:BB:BB:BB:BB", IPAddress: "192.168.1.2", State: "Offline"},
			{ID: "d-3", DeviceName: "Third", Description: "desc", MacAddress: "CC:CC:CC:CC:CC:CC", IPAddress: "192.168.1.3", State: "Offline"},
		},
		Groups: []config.Group{
			{ID: "g-1", GroupName: "MyGroup", Devices: []string{"d-1", "d-2", "d-3"}},
		},
	}
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(got.Groups))
	}
	wantDevices := []string{"d-1", "d-2", "d-3"}
	for i, want := range wantDevices {
		if got.Groups[0].Devices[i] != want {
			t.Errorf("group device[%d]: expected %q, got %q", i, want, got.Groups[0].Devices[i])
		}
	}
}

func TestUpdateDevice(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{
			{ID: "d-1", DeviceName: "OldName", Description: "desc", MacAddress: "AA:AA:AA:AA:AA:AA", IPAddress: "192.168.1.1", State: "Offline"},
		},
		Groups: []config.Group{},
	}
	config.WriteConfig(cfg)

	cfg.Devices[0].DeviceName = "NewName"
	cfg.Devices[0].IPAddress = "10.0.0.1"
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(got.Devices))
	}
	if got.Devices[0].DeviceName != "NewName" {
		t.Errorf("expected DeviceName %q, got %q", "NewName", got.Devices[0].DeviceName)
	}
	if got.Devices[0].IPAddress != "10.0.0.1" {
		t.Errorf("expected IPAddress %q, got %q", "10.0.0.1", got.Devices[0].IPAddress)
	}
}

func TestDeleteDevice(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{
			{ID: "d-1", DeviceName: "Keep", Description: "desc", MacAddress: "AA:AA:AA:AA:AA:AA", IPAddress: "192.168.1.1", State: "Offline"},
			{ID: "d-2", DeviceName: "Remove", Description: "desc", MacAddress: "BB:BB:BB:BB:BB:BB", IPAddress: "192.168.1.2", State: "Offline"},
		},
		Groups: []config.Group{},
	}
	config.WriteConfig(cfg)

	cfg.Devices = cfg.Devices[:1]
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Devices) != 1 {
		t.Fatalf("expected 1 device after deletion, got %d", len(got.Devices))
	}
	if got.Devices[0].DeviceName != "Keep" {
		t.Errorf("expected %q to remain, got %q", "Keep", got.Devices[0].DeviceName)
	}
}

func TestDeleteGroup(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{},
		Groups: []config.Group{
			{ID: "g-1", GroupName: "Keep", Devices: []string{}},
			{ID: "g-2", GroupName: "Remove", Devices: []string{}},
		},
	}
	config.WriteConfig(cfg)

	cfg.Groups = cfg.Groups[:1]
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Groups) != 1 {
		t.Fatalf("expected 1 group after deletion, got %d", len(got.Groups))
	}
	if got.Groups[0].GroupName != "Keep" {
		t.Errorf("expected %q to remain, got %q", "Keep", got.Groups[0].GroupName)
	}
}

func TestEmptyConfigReturnsEmptySlices(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	got := config.ReadConfig()

	if got.Devices == nil {
		t.Error("expected non-nil Devices slice on empty DB")
	}
	if got.Groups == nil {
		t.Error("expected non-nil Groups slice on empty DB")
	}
	if len(got.Devices) != 0 {
		t.Errorf("expected 0 devices, got %d", len(got.Devices))
	}
	if len(got.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(got.Groups))
	}
}

func TestGroupWithNoDevices(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{},
		Groups: []config.Group{
			{ID: "g-1", GroupName: "Empty", Devices: []string{}},
		},
	}
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if len(got.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(got.Groups))
	}
	if len(got.Groups[0].Devices) != 0 {
		t.Errorf("expected 0 devices in group, got %d", len(got.Groups[0].Devices))
	}
}

func TestDeviceStatePreserved(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	cfg := config.Config{
		Devices: []config.Device{
			{ID: "d-1", DeviceName: "Server", Description: "desc", MacAddress: "AA:AA:AA:AA:AA:AA", IPAddress: "192.168.1.1", State: "Online"},
		},
		Groups: []config.Group{},
	}
	config.WriteConfig(cfg)

	got := config.ReadConfig()

	if got.Devices[0].State != "Online" {
		t.Errorf("expected State %q, got %q", "Online", got.Devices[0].State)
	}
}
