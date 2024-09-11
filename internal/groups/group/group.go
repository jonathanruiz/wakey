package group

import (
	"fmt"
	"strings"
	"wakey/internal/common/status"
	"wakey/internal/common/style"
	"wakey/internal/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

var (
	focusedButton = style.FocusedStyle.Render("[ Submit ]")                    // The focused button
	blurredButton = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Submit")) // The blurred button
)

// Model is the model for the Group component
type Model struct {
	focusIndex    int
	inputs        []textinput.Model
	err           []error
	previousModel tea.Model
	currentConfig config.Config
	keys          keyMap
	help          help.Model
	selectedRow   []string
	deviceNameMap map[string]string
	deviceIDMap   map[string]string
}

// InitialModel returns the initial model for the Group component
func InitialModel(previousModel tea.Model, selectedRow ...[]string) Model {
	m := Model{
		err:           make([]error, 2),           // Initialize the slice with length 2
		inputs:        make([]textinput.Model, 2), // Initialize the slice with length 2
		currentConfig: config.ReadConfig(),
		keys:          keys,
		help:          help.New(),
		previousModel: previousModel,
		deviceNameMap: createDeviceNameMap(config.ReadConfig().Devices), // Initialize the deviceNameMap
		deviceIDMap:   createDeviceIDMap(config.ReadConfig().Devices),   // Initialize the deviceIDMap
	}

	var deviceNames []string

	// Check if this is an edit operation
	if len(selectedRow) > 0 {
		m.selectedRow = selectedRow[0]

		// Get the device names
		deviceIDs := strings.Split(selectedRow[0][2], ",")

		// Remove any leading or trailing spaces
		for i, group := range deviceIDs {
			deviceIDs[i] = strings.TrimSpace(group)
		}

		// Convert the device IDs to device names
		deviceNames = convertDeviceIDsToNames(deviceIDs, m.deviceIDMap)

		// Set the device names in the input field
		selectedRow[0][2] = strings.Join(deviceNames, ",")
	}

	// Create a new text input model for each input field
	var ti textinput.Model

	// Loop through the inputs and create a new text input model for each
	for i := range m.inputs {
		ti = textinput.New()
		ti.Cursor.Style = style.FocusedStyle
		ti.CharLimit = 64

		switch i {
		// Group name
		case 0:
			ti.Prompt = "Group Name   : "
			ti.Placeholder = "Enter the group name"
			ti.Focus()
			ti.PromptStyle = style.FocusedStyle
			ti.TextStyle = style.FocusedStyle

			if selectedRow != nil {
				ti.SetValue(selectedRow[0][1])
			}
		// Devices
		case 1:
			ti.Prompt = "Devices      : "
			ti.Placeholder = "Device1, Device2"
			ti.ShowSuggestions = true
			ti.SetSuggestions(deviceNames)

			if selectedRow != nil {
				devices := strings.Split(selectedRow[0][2], ",")

				ti.SetValue(strings.Join(devices, ", "))
			}
		}

		// Add the textinput model to the slice
		m.inputs[i] = ti
	}

	return m
}

// Update function for the Group model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update function for the Group model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Return to the list
		case key.Matches(msg, m.keys.Quit):
			return m.previousModel, nil

		// Toggle help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// Set focus to next input
		case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.Enter):
			// Check if the user pressed enter with the submit button focused
			if key.Matches(msg, m.keys.Enter) && m.focusIndex == len(m.inputs) {

				// Run the validators
				m.err[0] = m.groupNameValidator(m.inputs[0].Value())
				m.err[1] = m.devicesValidator(m.inputs[1].Value())

				if m.focusIndex == len(m.inputs) {
					// Handle form submission
					// Reset focus index or update state as needed
					m.focusIndex = 0

					// Validate the group name
					if !m.validateInput(0, m.groupNameValidator) {
						return m, nil
					}

					// Validate the devices
					if !m.validateInput(1, m.devicesValidator) {
						return m, nil
					}

					// Load existing devices
					existingDevices := createDeviceIDMap(m.currentConfig.Devices)

					// Convert the Group value from string to []string
					deviceValue := strings.Split(m.inputs[1].Value(), ",")

					// Remove any leading or trailing spaces
					for i, group := range deviceValue {
						deviceValue[i] = strings.TrimSpace(group)
					}

					// Replace the device names with device IDs
					deviceValue = convertDeviceNamesToIDs(deviceValue, existingDevices)

					// Check if we are editing an existing group
					if m.selectedRow != nil {
						// Get the selected group
						selected := m.selectedRow

						// Update the group in the config
						for i, group := range m.currentConfig.Groups {
							if group.ID == selected[0] {
								m.currentConfig.Groups[i] = config.Group{
									ID:        group.ID,
									GroupName: m.inputs[0].Value(),
									Devices:   deviceValue,
								}
								break
							}
						}
					} else {
						// Append the group to the config
						updatedGroups := append(m.currentConfig.Groups, config.Group{
							ID:        uuid.NewString(),
							GroupName: m.inputs[0].Value(),
							Devices:   deviceValue,
						})

						// Create a new config with the updated group
						m.currentConfig = config.Config{
							Devices: m.currentConfig.Devices,
							Groups:  updatedGroups,
						}
					}

					// Write the the new version of the config to the file
					config.WriteConfig(m.currentConfig)

					// Set the status message
					status.Message = fmt.Errorf("group [%s] added", m.inputs[0].Value())

					// Return to the list and clear the screen
					return m.previousModel, func() tea.Msg {
						return tea.ClearScreen()
					}
				}
			}

			// Cycle indexes
			if key.Matches(msg, m.keys.Up) {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// Wrap around
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			// Set focus to the input
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = style.FocusedStyle
					m.inputs[i].TextStyle = style.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = style.NoStyle
				m.inputs[i].TextStyle = style.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

// Define the DeleteGroupPopup function
func DeleteGroupPopup(groupName, macAddress string, m tea.Model) (tea.Model, tea.Cmd) {
	// Create a popup message for confirmation
	popupMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render("Are you sure you want to delete " + groupName + "?")

	// Return the popup message and a command to handle user input
	return m, func() tea.Msg {
		return popupMessage
	}
}

// updateInputs updates all the text inputs in the Group model.
func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Ensure focus is correctly managed
	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex].Focus()
	}

	return tea.Batch(cmds...)
}

// View function for the Group model
func (m Model) View() string {
	// The header
	s := "\n"

	buttons := style.FocusedTab.Render("Groups > New Group")

	s = lipgloss.PlaceHorizontal(style.TermWidth, lipgloss.Center, buttons) + "\n"

	// Render the inputs
	for i, input := range m.inputs {
		// Check if there are any errors and if the errors are not nil
		if len(m.err) > 0 && m.err[i] != nil {
			// Display the error message inline with the first input field
			// Refer to discussion: https://github.com/charmbracelet/bubbles/discussions/306
			s += lipgloss.JoinHorizontal(lipgloss.Left, input.View()+"   ", style.ErrStyle(m.err[i].Error())) + "\n"
		} else {
			// Display the input field
			s += input.View() + "\n"
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	s += fmt.Sprintf("\n\n%s\n\n", *button)

	// Render the help text
	s += m.help.View(m.keys)

	return s

}

// createDeviceNameMap creates a map of device IDs to device names
func createDeviceNameMap(devices []config.Device) map[string]string {
	deviceNameMap := make(map[string]string)
	for _, device := range devices {
		deviceNameMap[device.DeviceName] = device.ID
	}
	return deviceNameMap
}

func createDeviceIDMap(devices []config.Device) map[string]string {
	deviceIDMap := make(map[string]string)
	for _, device := range devices {
		deviceIDMap[device.ID] = device.DeviceName
	}
	return deviceIDMap
}

// convertDeviceNamesToIDs converts a slice of device names to a slice of device IDs
func convertDeviceNamesToIDs(deviceNames []string, deviceNameMap map[string]string) []string {
	var deviceIDs []string
	for _, name := range deviceNames {
		for id, deviceName := range deviceNameMap {
			if deviceName == name {
				deviceIDs = append(deviceIDs, id)
				break
			}
		}
	}
	return deviceIDs
}

// convertDeviceIDsToNames converts a slice of device IDs to a slice of device names
func convertDeviceIDsToNames(deviceIDs []string, deviceNameMap map[string]string) []string {
	var deviceNames []string
	for _, id := range deviceIDs {
		if name, ok := deviceNameMap[id]; ok {
			deviceNames = append(deviceNames, name)
		}
	}
	return deviceNames
}
