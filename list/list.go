package list

import (
	"fmt"
	"strconv"
	"wakey/config"
	"wakey/device"
	"wakey/popup"
	"wakey/style"
	"wakey/wol"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the Device component
type Model struct {
	devices []config.Device // list of devices to wake
	keys    keyMap
	help    help.Model
	table   table.Model
	status  error
}

// InitialModel function for the Device model
func InitialModel() tea.Model {
	// Create the config file if it	doesn't exist
	status := config.CreateConfig()

	// Update the Status of the devices
	devices := config.UpdateStatus().Devices

	// Define table columns
	columns := []table.Column{
		{Title: "Device", Width: 20},
		{Title: "Description", Width: 30},
		{Title: "MAC Address", Width: 20},
		{Title: "IP Address", Width: 15},
		{Title: "Status", Width: 15},
	}

	// Define table rows
	rows := make([]table.Row, len(devices))
	for i, device := range devices {
		rows[i] = table.Row{
			device.DeviceName,
			device.Description,
			device.MacAddress,
			device.IPAddress,
			device.Status,
		}
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Set the custom key bindings
	t.KeyMap = table.KeyMap{
		LineUp:   keys.Up,
		LineDown: keys.Down,
	}

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
		keys:   keys,
		help:   help.New(),
		table:  t,
		status: status,
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
			device.Status,
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

		// Delete device
		case key.Matches(msg, m.keys.Delete):
			// Get the selected device
			selected := m.table.SelectedRow()

			// Return popup message for confirmation
			return popup.NewPopupMsg("Are you sure you want to delete "+selected[0]+" ("+selected[2]+")?", m, m.table), nil

		// Refresh the table
		case key.Matches(msg, m.keys.Refresh):
			// clear screen
			return InitialModel(), nil

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

			// Write the status message
			deviceName := selected[0]
			m.status = fmt.Errorf("Waking up " + deviceName + " (" + macAddress + ")")

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
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
	for _, device := range newConfig.Devices {
		// Append the device to the rows
		// This will make sure to output all the data for the device
		// The order of the columns must match the order of the columns in the table
		rows = append(rows, table.Row{device.DeviceName, device.Description, device.MacAddress, device.IPAddress, device.Status})
	}

	// Update the table with the new rows
	m.table.SetRows(rows)

	// The header
	s := style.TitleStyle.Render("Which device should you wake?") + "\n\n"

	// Render the table
	s += m.table.View() + "\n"

	// Show device count
	s += style.DeviceCountStyle.Render(" Number of devices: "+strconv.Itoa(len(m.table.Rows()))) + "\n" // srtconv.Itoa converts int to string

	// Status message
	var statusMessage string
	if m.status != nil {
		statusMessage = m.status.Error()
	} else {
		statusMessage = "No status"
	}
	s += style.StatusStyle.Render("Status: "+style.StatusMessageStyle.Render(statusMessage)) + "\n"

	// Help text
	s += m.help.View(m.keys)

	return s
}
