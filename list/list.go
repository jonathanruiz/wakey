package list

import (
	"fmt"
	"strconv"
	"wakey/config"
	"wakey/device"
	"wakey/style"
	"wakey/wol"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

// Keybindings for the Device component
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
		key.WithHelp("enter", "select"),
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

// Model for the Device component
type Model struct {
	devices []config.Device // list of devices to wake
	keys    keyMap
	help    help.Model
	table   table.Model
}

// InitialModel function for the Device model
func InitialModel() tea.Model {
	// Get devices from config
	devices := config.ReadConfig().Devices

	// Define table columns
	columns := []table.Column{
		{Title: "Device", Width: 20},
		{Title: "Description", Width: 30},
		{Title: "MAC Address", Width: 20},
		{Title: "IP Address", Width: 15},
	}

	// Define table rows
	rows := make([]table.Row, len(devices))
	for i, device := range devices {
		rows[i] = table.Row{
			device.DeviceName,
			device.Description,
			device.MacAddress,
			device.IPAddress,
		}
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Get the default table styles
	s := style.DefaultTableStyles()

	// Set the styles
	t.SetStyles(table.Styles{
		Header:   s.Header,
		Selected: s.Selected,
	})

	return Model{
		// A list of devices to wake. This could be fetched from a database or config file
		devices: devices,
		// A map which indicates which devices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `devices` slice, above.
		keys:  keys,
		help:  help.New(),
		table: t,
	}
}

// Init function for the Device model
func (m Model) Init() tea.Cmd { return nil }

// Update function for the Device model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Get new number of rows
	rows := make([]table.Row, len(config.ReadConfig().Devices))

	// Update the table with the new rows
	m.table.SetRows(rows)

	// Define table rows
	for i, device := range config.ReadConfig().Devices {
		m.table.Rows()[i] = table.Row{
			device.DeviceName,
			device.Description,
			device.MacAddress,
			device.IPAddress,
		}
	}

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

		// Wake device
		case key.Matches(msg, m.keys.Enter):
			// Get the selected device
			selected := m.table.SelectedRow()

			// Get the device MAC address
			macAddress := selected[2]

			// Wake the device
			wol.WakeDevice(macAddress)

			// Show modal with device name
			deviceName := selected[0]
			popupModel := NewPopupMsg(fmt.Sprintf("Magic Packet has been sent to: %s", deviceName), m)

			return popupModel, popupModel.Init()
		}
	}

	// Update the table
	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View function for the Device model
func (m Model) View() string {
	// Get updated config file
	newConfig := config.ReadConfig()

	// Convert m.devices from []string to []table.Row
	var rows []table.Row
	for _, choice := range newConfig.Devices {
		// Append the device to the rows
		// This will make sure to output all the data for the device
		// The order of the columns must match the order of the columns in the table
		rows = append(rows, table.Row{choice.DeviceName, choice.Description, choice.MacAddress, choice.IPAddress})
	}

	// Update the table with the new rows
	m.table.SetRows(rows)

	// The header
	s := style.TitleStyle.Render("Which device should you wake?") + "\n\n"

	// Render the table
	s += m.table.View() + "\n"

	// Show device count
	s += style.DeviceCountStyle.Render("Number of devices: "+strconv.Itoa(len(m.table.Rows()))) + "\n\n" // srtconv.Itoa converts int to string

	// Help text
	s += m.help.View(m.keys)

	return s
}
