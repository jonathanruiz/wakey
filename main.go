package main

import (
	"fmt"
	"os"

	"wakey/config"
	"wakey/list"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the config file if it	doesn't exist
	config.CreateConfig()

	// Create a new program and open the alternate screen
	p := tea.NewProgram(list.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
