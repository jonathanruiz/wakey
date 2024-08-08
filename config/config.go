package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Devices []string `json:"devices"`
}

var (
	HomeDir, HomeDirErr = os.UserHomeDir()
	ConfigPath          = filepath.Join(HomeDir, ".wakey_config.json")
)

// Create a config file if it doesn't exist in the users home directory.
// Returns the contents of the config file.
func CreateConfig() {

	// Check if we got an error
	if HomeDirErr != nil {
		fmt.Println("Error getting home directory:", HomeDirErr)
		return
	}

	// Create the path to the config file
	configPath := filepath.Join(HomeDir, ".wakey_config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If it doesn't exist, create it
		config := Config{
			Devices: []string{},
		}

		// Marshal the config to JSON
		data, err := json.MarshalIndent(config, "", "  ")

		// Check if we got an error
		if err != nil {
			fmt.Println("Error marshalling config:", err)
			return
		}

		// Write the config to the file
		err = os.WriteFile(configPath, data, 0644)

		// Check if we got an error
		if err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}

		// Print a message to the user
		fmt.Println("Config file created at", configPath)
	} else {
		// Print a message to the user
		fmt.Println("Config file already exists at", configPath)
	}
}

// Read the config file and return the contents.
func ReadConfig() Config {
	// Read the config file
	data, err := os.ReadFile(ConfigPath)

	// Check if we got an error
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return Config{}
	}

	// Unmarshal the config
	var config Config
	err = json.Unmarshal(data, &config)

	// Check if we got an error
	if err != nil {
		fmt.Println("Error unmarshalling config:", err)
		return Config{}
	}

	// Return the config
	return config
}

/*
Convert the config to a JSON string.

Helpful for debugging and seeing the contents of the config.
Might not be necessary anymore since `ReadConfig()` returns a `Config` struct as JSON.
*/
func (c Config) ConfigToString() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config:", err)
		return ""
	}
	return string(data)
}
