package main

import (
	"fmt"
	"os"

	"wakey/internal"
	"wakey/internal/common/status"
	"wakey/internal/config"

	tea "charm.land/bubbletea/v2"
)

func main() {
	status.Message = config.CreateConfig()
	// Create a new program and open the alternate screen
	p := tea.NewProgram(internal.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
