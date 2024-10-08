package device

import (
	"fmt"
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

// Model is the model for the Device component
type Model struct {
	focusIndex    int
	inputs        []textinput.Model
	err           []error
	previousModel tea.Model
	currentConfig config.Config
	keys          keyMap
	help          help.Model
	selectedRow   []string
}

// InitialModel returns the initial model for the Device component
func InitialModel(previousModel tea.Model, selectedRow ...[]string) Model {
	m := Model{
		err:           make([]error, 4),           // Initialize the slice with length 4
		inputs:        make([]textinput.Model, 4), // Initialize the slice with length 4
		currentConfig: config.ReadConfig(),
		keys:          keys,
		help:          help.New(),
		previousModel: previousModel,
	}

	// Check if this is an edit operation
	if len(selectedRow) > 0 {
		// Set the selected row
		m.selectedRow = selectedRow[0]
	}

	// Create a new text input model for each input field
	var ti textinput.Model

	// Loop through the inputs and create a new text input model for each
	for i := range m.inputs {
		ti = textinput.New()
		ti.Cursor.Style = style.FocusedStyle
		ti.CharLimit = 64

		switch i {
		// Device name
		case 0:
			ti.Prompt = "Device Name   : "
			ti.Placeholder = "Enter the device name"
			ti.Focus()
			ti.PromptStyle = style.FocusedStyle
			ti.TextStyle = style.FocusedStyle

			if selectedRow != nil {
				ti.SetValue(selectedRow[0][1])
			}
		// Description
		case 1:
			ti.Prompt = "Description   : "
			ti.Placeholder = "Enter a description for the device"

			if selectedRow != nil {
				ti.SetValue(selectedRow[0][2])
			}
		// MAC address
		case 2:
			ti.Prompt = "MAC Address   : "
			ti.Placeholder = "00:00:00:00:00:00"

			if selectedRow != nil {
				ti.SetValue(selectedRow[0][3])
			}
		// IP address
		case 3:
			ti.Prompt = "IP Address    : "
			ti.Placeholder = "0.0.0.0"

			if selectedRow != nil {
				ti.SetValue(selectedRow[0][4])
			}
		}

		// Add the textinput model to the slice
		m.inputs[i] = ti
	}

	return m
}

// Update function for the Device model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update function for the Device model
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
				m.err[0] = m.deviceNameValidator(m.inputs[0].Value())
				m.err[1] = m.descriptionValidator(m.inputs[1].Value())
				m.err[2] = m.macAddressValidator(m.inputs[2].Value())
				m.err[3] = m.ipAddressValidator(m.inputs[3].Value())

				if m.focusIndex == len(m.inputs) {
					// Handle form submission
					// Reset focus index or update state as needed
					m.focusIndex = 0

					// Validate the device name
					if !m.validateInput(0, m.deviceNameValidator) {
						return m, nil
					}

					if !m.validateInput(1, m.descriptionValidator) {
						return m, nil
					}

					if !m.validateInput(2, m.macAddressValidator) {
						return m, nil
					}

					if !m.validateInput(3, m.ipAddressValidator) {
						return m, nil
					}

					// Check if we are editing an existing device
					if m.selectedRow != nil {
						// Get the selected device
						selected := m.selectedRow

						// Update the device in the config
						for i, device := range m.currentConfig.Devices {
							if device.ID == selected[0] {
								m.currentConfig.Devices[i] = config.Device{
									ID:          m.currentConfig.Devices[i].ID,
									DeviceName:  m.inputs[0].Value(),
									Description: m.inputs[1].Value(),
									MacAddress:  m.inputs[2].Value(),
									IPAddress:   m.inputs[3].Value(),
									State:       m.currentConfig.Devices[i].State,
								}
								break
							}
						}
					} else {
						// Append the device to the config
						updatedDevices := append(m.currentConfig.Devices, config.Device{
							ID:          uuid.NewString(),
							DeviceName:  m.inputs[0].Value(),
							Description: m.inputs[1].Value(),
							MacAddress:  m.inputs[2].Value(),
							IPAddress:   m.inputs[3].Value(),
							State:       "Offline",
						})

						// Create a new config with the updated devices
						m.currentConfig = config.Config{
							Devices: updatedDevices,
							Groups:  m.currentConfig.Groups,
						}
					}

					// Write the the new version of the config to the file
					config.WriteConfig(m.currentConfig)

					// Set the status message
					status.Message = fmt.Errorf("device [%s] (%s) added", m.inputs[0].Value(), m.inputs[2].Value())

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

// Define the DeleteDevicePopup function
func DeleteDevicePopup(deviceName, macAddress string, m tea.Model) (tea.Model, tea.Cmd) {
	// Create a popup message for confirmation
	popupMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render("Are you sure you want to delete " + deviceName + " (" + macAddress + ")?")

	// Return the popup message and a command to handle user input
	return m, func() tea.Msg {
		return popupMessage
	}
}

// updateInputs updates all the text inputs in the Device model.
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

// View function for the Device model
func (m Model) View() string {
	// The header
	s := "\n"

	buttons := style.FocusedTab.Render("Devices > New Device")

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

// CreateGroupNameMap creates a map of group names to group IDs
func CreateGroupNameMap(groups []config.Group) map[string]string {
	groupNameMap := make(map[string]string)
	for _, group := range groups {
		groupNameMap[group.ID] = group.GroupName
	}
	return groupNameMap
}

// createGroupIDMap creates a map of group IDs to group names
func createGroupIDMap(groups []config.Group) map[string]string {
	groupNameMap := make(map[string]string)
	for _, group := range groups {
		groupNameMap[group.GroupName] = group.ID
	}
	return groupNameMap
}
