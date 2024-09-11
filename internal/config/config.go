package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"wakey/internal/wol"
)

// Config struct for the config file.
type Device struct {
	ID          string `json:"ID"`
	DeviceName  string `json:"DeviceName"`
	Description string `json:"Description"`
	MacAddress  string `json:"MacAddress"`
	IPAddress   string `json:"IPAddress"`
	State       string `json:"State"`
}

type Group struct {
	ID        string   `json:"ID"`
	GroupName string   `json:"GroupName"`
	Devices   []string `json:"Devices"` // contains IDs of devices
}

// Config struct for the config file.
type Config struct {
	Devices []Device `json:"devices"`
	Groups  []Group  `json:"groups"`
}

var (
	HomeDir, HomeDirErr = os.UserHomeDir()                             // Get the users home directory
	ConfigPath          = filepath.Join(HomeDir, ".wakey_config.json") // Create the path to the config file
)

// Create a config file if it doesn't exist in the users home directory.
// Returns the contents of the config file.
func CreateConfig() error {

	// Check if we got an error
	if HomeDirErr != nil {

		return fmt.Errorf("error getting home directory: %v", HomeDirErr)
	}

	// Create the path to the config file
	configPath := filepath.Join(HomeDir, ".wakey_config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If it doesn't exist, create it
		config := Config{
			Devices: []Device{},
			Groups:  []Group{},
		}

		// Marshal the config to JSON
		data, err := json.MarshalIndent(config, "", "  ")

		// Check if we got an error
		if err != nil {
			return fmt.Errorf("error marshalling config: %v", err)
		}

		// Write the config to the file
		err = os.WriteFile(configPath, data, 0644)

		// Check if we got an error
		if err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		// Print a message to the user
		return fmt.Errorf("Config file created at: %v", configPath)
	} else {
		// Print a message to the user
		return fmt.Errorf("Config file already exists at: %v", configPath)
	}
}

// Read the config file and return the contents.
func ReadConfig() Config {
	// Check if we got an error
	if HomeDirErr != nil {
		fmt.Println("Error getting home directory:", HomeDirErr)
		return Config{}
	}

	// Read the config file
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		fmt.Println("error reading config file:", err)
		return Config{}
	}

	// Unmarshal the JSON data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config:", err)
		return Config{}
	}

	return config
}

// Write the config to the config file.
func WriteConfig(config Config) {
	// Marshal the config to JSON
	data, err := json.MarshalIndent(config, "", "  ")

	// Check if we got an error
	if err != nil {
		fmt.Println("Error marshalling config:", err)
		return
	}

	// Write the config to the file
	err = os.WriteFile(ConfigPath, data, 0644)

	// Check if we got an error
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

	// Print a message to the user
	fmt.Println("Config file updated at", ConfigPath)
}

// Update the State of the devices
func GetUpdateState() Config {
	// Get the devices
	cfg := ReadConfig()
	devices := cfg.Devices
	groups := cfg.Groups

	// Loop through the devices
	for i, device := range devices {
		// Get the State of the device
		isOnline := wol.IsOnline(device.IPAddress)

		// Update the State of the device
		if isOnline {
			devices[i].State = "Online"
		} else {
			devices[i].State = "Offline"
		}
	}

	// Write the updated config file
	WriteConfig(Config{Devices: devices, Groups: groups})

	// Return the config file
	return Config{Devices: devices, Groups: groups}
}

/*
Convert the config to a JSON string.

Helpful for debugging and seeing the contents of the config.
Might not be necessary anymore since `ReadConfig()` returns a `Config` struct as JSON.
*/
func (c Config) ConfigToString() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println("error marshalling config:", err)
		return ""
	}
	return string(data)
}
