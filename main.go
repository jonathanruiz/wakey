package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"wakey/list"

	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Devices []string `json:"devices"`
}

// Create a config file if it doesn't exist in the users home directory.
// return the contents of the config file.
func CreateConfig() {
	// Get the users home directory
	homeDir, err := os.UserHomeDir()

	// Check if we got an error
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create the path to the config file
	configPath := filepath.Join(homeDir, ".wakey_config.json")

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
		// If it exists, read it
		data, err := os.ReadFile(configPath)

		// Check if we got an error
		if err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		// Unmarshal the config
		var config Config
		err = json.Unmarshal(data, &config)

		// Check if we got an error
		if err != nil {
			fmt.Println("Error unmarshalling config:", err)
			return
		}

		// Print a message to the user
		fmt.Println("Config file loaded with devices:", config.Devices)

	}

}

func main() {
	homeDir, err := os.UserHomeDir()

	// Create the config file
	CreateConfig()

	configPath := filepath.Join(homeDir, ".wakey_config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	println(string(data))

	p := tea.NewProgram(list.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
