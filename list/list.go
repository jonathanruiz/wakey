package list

import (
	"wakey/config"
	"wakey/device"
	"wakey/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
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
		{k.Up, k.Down},   // first column
		{k.Enter, k.New}, // second column
		{k.Help, k.Quit}, // third column
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
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "toggle select"),
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
	choices    []string // list of devices to wake
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	table      table.Model
}

func InitialModel() tea.Model {
	// Get devices from config
	choices := config.ReadConfig().Devices

	// Define table columns
	columns := []table.Column{
		{Title: "Device", Width: 20},
		{Title: "Description", Width: 30},
		{Title: "MAC Address", Width: 20},
		{Title: "IP Address", Width: 15},
	}

	// Define table rows
	rows := make([]table.Row, len(choices))
	for i, device := range choices {
		rows[i] = table.Row{
			device,
		}
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return Model{
		// A list of devices to wake. This could be fetched from a database or config file
		choices: choices,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		keys:  keys,
		help:  help.New(),
		table: t,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now".
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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

		// Toggle help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	// Update the table
	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Get updated config file
	newConfig := config.ReadConfig()

	// Convert m.choices from []string to []table.Row
	var rows []table.Row
	for _, choice := range newConfig.Devices {
		rows = append(rows, table.Row{choice})
	}

	// Update the table with the new rows
	m.table.SetRows(rows)

	// The header
	s := style.TitleStyle.Render("Which device should you wake?") + "\n\n"

	// Render the table
	s += m.table.View() + "\n"

	// Help text
	s += m.help.View(m.keys)

	return s
}
