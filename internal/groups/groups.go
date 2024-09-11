package groups

import (
	"fmt"
	"strings"
	"wakey/internal/config"
	"wakey/internal/groups/group"
	"wakey/internal/helper/popup"
	"wakey/internal/helper/status"
	"wakey/internal/helper/style"
	"wakey/internal/helper/wol"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the Group component
type Model struct {
	groups        []config.Group // list of groups
	previousModel tea.Model
	keys          keyMap
	help          help.Model
	table         table.Model
}

// Init function for the Device model
func (m Model) Init() tea.Cmd { return nil }

// InitialModel function for the Group model
func InitialModel(previousModel tea.Model) tea.Model {
	// Get groups with updated state
	groups := config.GetUpdateState().Groups

	// Define table columns
	columns := []table.Column{
		{Title: "ID", Width: 0},
		{Title: "Group Name", Width: 20},
		{Title: "Devices", Width: 30},
	}

	// Define table rows
	rows := make([]table.Row, len(groups))
	for i, group := range groups {
		// Create a string of devices in the group
		devicesInGroup := ""
		for _, device := range group.Devices {
			devicesInGroup += device + ", "
		}

		// Remove the trailing comma and space if devicesInGroup is not empty
		if len(devicesInGroup) > 0 {
			devicesInGroup = devicesInGroup[:len(devicesInGroup)-2]
		}

		rows[i] = table.Row{
			group.GroupName,
			devicesInGroup,
		}
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
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
		groups: groups,
		// A map which indicates which devices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `devices` slice, above.
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
		table:         t,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Get new number of rows
	rows := make([]table.Row, len(config.ReadConfig().Groups))

	// Update the table with the new rows
	m.table.SetRows(rows)

	// Define table rows
	for i, group := range config.ReadConfig().Groups {
		deviceValue := strings.Join(group.Devices, ", ")
		m.table.Rows()[i] = table.Row{
			group.ID,
			group.GroupName,
			deviceValue,
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Create):
			// Create a new group
			return group.InitialModel(m), nil
		case key.Matches(msg, m.keys.Edit):
			// Edit the selected group
			selected := m.table.SelectedRow()

			return group.InitialModel(m, selected), nil

		case key.Matches(msg, m.keys.Delete):
			// Delete the selected group
			selected := m.table.SelectedRow()

			// Return popup message for confirmation
			return popup.NewPopupMsg("Are you sure you want to delete "+selected[1]+"?", m, m.table, deleteGroup), nil

		case key.Matches(msg, m.keys.Enter):
			// Extract the selected group and get the device IDs
			selected := m.table.SelectedRow()
			deviceIDs := selected[2]

			// Split the device IDs into an array
			deviceIDsArr := strings.Split(deviceIDs, ", ")

			// Create a map of device IDs to MAC addresses
			deviceMap := createDeviceMacAddressMap(config.ReadConfig().Devices)

			// Get the MAC addresses for the device IDs
			macAddresses := getMacAddresses(deviceIDsArr, deviceMap)

			// Wake the group using the MAC addresses
			err := wol.WakeGroup(macAddresses)
			if err != nil {
				status.Message = err
			} else {
				status.Message = fmt.Errorf("waking [%s] group", selected[1])
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.View):
			// Switch to the device view
			return m.previousModel, tea.ClearScreen

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Refactored code from lines 171 to 195
	const maxRows = 10 // Define the maximum number of rows to display

	// Get updated config file
	newConfig := config.ReadConfig()

	// Create a map of device IDs to device names
	deviceNameMap := createDeviceNameMap(newConfig.Devices)

	var rows []table.Row
	for _, group := range newConfig.Groups {
		var deviceNames []string
		for _, deviceID := range group.Devices {
			if deviceName, ok := deviceNameMap[deviceID]; ok {
				deviceNames = append(deviceNames, deviceName)
			} else {
				deviceNames = append(deviceNames, deviceID) // Fallback to device ID if name not found
			}
		}
		deviceNamesStr := strings.Join(deviceNames, ", ")

		rows = append(rows, table.Row{group.ID, group.GroupName, deviceNamesStr})
	}

	// Truncate rows if they exceed the maximum number
	if len(rows) > maxRows {
		rows = rows[:maxRows]
	}

	// Update the table with the new rows
	m.table.SetRows(rows)

	// The header
	s := "\n"

	// Render the table
	s += m.table.View() + "\n"

	// Status message
	var statusMessage string
	if status.Message != nil {
		statusMessage = status.Message.Error()
	} else {
		statusMessage = "No status"
	}
	s += style.StatusStyle.Render("Status: "+style.StatusMessageStyle.Render(statusMessage)) + "\n"

	// Help text
	s += m.help.View(m.keys)
	return s
}

// createDeviceAttributeMap creates a map of device IDs to a specified attribute
func createDeviceAttributeMap(devices []config.Device, attributeFunc func(config.Device) string) map[string]string {
	deviceMap := make(map[string]string)
	for _, device := range devices {
		deviceMap[device.ID] = attributeFunc(device)
	}
	return deviceMap
}

// createDeviceNameMap creates a map of device IDs to device names
func createDeviceNameMap(devices []config.Device) map[string]string {
	return createDeviceAttributeMap(devices, func(d config.Device) string {
		return d.DeviceName
	})
}

// createDeviceMacAddressMap creates a map of device IDs to MAC addresses
func createDeviceMacAddressMap(devices []config.Device) map[string]string {
	return createDeviceAttributeMap(devices, func(d config.Device) string {
		return d.MacAddress
	})
}

// getMacAddresses returns the MAC addresses for the given device IDs
func getMacAddresses(deviceIDs []string, deviceMap map[string]string) []string {
	var macAddresses []string
	for _, deviceID := range deviceIDs {
		if macAddress, ok := deviceMap[deviceID]; ok {
			macAddresses = append(macAddresses, macAddress)
		}
	}
	return macAddresses
}

func deleteGroup(selectedRow []string) (string, error) {
	currentConfig := config.ReadConfig()
	for i, group := range currentConfig.Groups {
		if group.ID == selectedRow[0] {
			currentConfig.Groups = append(currentConfig.Groups[:i], currentConfig.Groups[i+1:]...)
			config.WriteConfig(currentConfig)
			return fmt.Sprintf("group [%s] removed", selectedRow[1]), nil
		}
	}
	return "", fmt.Errorf("group [%s] not found", selectedRow[1])
}
