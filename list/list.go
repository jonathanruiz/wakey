package list

import (
	"fmt"
	"wakey/config"
	"wakey/device"
	"wakey/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	New   key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.New, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Enter, k.New},                // second column
		{k.Help, k.Quit},                // third column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "toggle select"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new device"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Model struct {
	choices    []string         // list of devices to wake
	cursor     int              // which device is selected
	selected   map[int]struct{} // which devices are selected
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func InitialModel() tea.Model {
	// Get devices from config
	choices := config.ReadConfig().Devices

	return Model{
		// A list of devices to wake. This could be fetched from a database or config file
		choices: choices,
		cursor:  0,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected:   make(map[int]struct{}),
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now".
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Get devices from config
	m.choices = config.ReadConfig().Devices

	switch msg := msg.(type) {

	// Check if it was a key press
	case tea.KeyMsg:

		// Check which key was pressed
		switch {

		// Create new device
		case key.Matches(msg, m.keys.New):
			return device.InitialModel(func() tea.Model { return m }), nil

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar toggle
		// the selected state for the item that the cursor is pointing at.
		case key.Matches(msg, m.keys.Enter):
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		// Toggle help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	// Return the updated Model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// Get updated config file
	newConfig := config.ReadConfig()

	// Update the choices with the new config
	m.choices = newConfig.Devices

	// The header
	s := style.TitleStyle.Render("Which device should you wake?") + "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = style.CursorStyle.Render(">")
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = style.CursorStyle.Render("x")
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// Help text
	s += m.help.View(m.keys)

	return s
}
