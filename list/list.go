package list

import (
	"fmt"
	"wakey/config"
	"wakey/device"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	viewport viewport.Model
	choices  []string         // list of devices to wake
	cursor   int              // which device is selected
	selected map[int]struct{} // which devices are selected
}

func InitialModel(config config.Config) tea.Model {

	// Create a new viewport
	vp := viewport.New(20, 10) // Adjust width and height as needed

	// Get devices from config
	choices := config.Devices

	return Model{
		viewport: vp,
		// A list of devices to wake. This could be fetched from a database or config file
		choices: choices,
		cursor:  0,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m *Model) AddChoice(choice string) {
	m.choices = append(m.choices, choice)
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

		// Create new device
		case "n":
			return device.InitialModel(func() tea.Model { return m }, m.AddChoice), nil

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated Model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header
	s := "Which device should you wake?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress n to add new device."
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	s += "\n" + m.viewport.View()
	return s
}
