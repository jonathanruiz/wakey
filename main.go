package main

import (
	"fmt"
	"os"

	"wakey/list"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create a new program and open the alternate screen
	p := tea.NewProgram(list.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
