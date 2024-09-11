package main

import (
	"fmt"
	"os"

	"wakey/internal"
	"wakey/internal/common/status"
	"wakey/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	status.Message = config.CreateConfig()
	// Create a new program and open the alternate screen
	p := tea.NewProgram(internal.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
