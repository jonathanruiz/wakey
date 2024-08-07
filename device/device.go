package device

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	name        string
	description string
	mac_address string
	ip_address  string
}

func InitialModel() tea.Model {
	return Model{
		name:        "",
		description: "",
		mac_address: "",
		ip_address:  "",
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now".
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Check if it was a key press
	case tea.KeyMsg:

		// Check which key was pressed
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// Return the updated Model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header
	s := "New Device\n\n"
	s += "\nPress q to quit.\n"

	return s
}
