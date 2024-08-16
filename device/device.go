package device

import (
	"fmt"
	"regexp"
	"wakey/config"
	"wakey/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedButton   = style.FocusedStyle.Render("[ Submit ]")                         // The focused button
	blurredButton   = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Submit"))      // The blurred button
	macAddressRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`) // The regex for the MAC address
	ipAddressRegex  = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)                   // The regex for the IP address
)

// Model is the model for the Device component
type Model struct {
	focusIndex    int
	inputs        []textinput.Model
	err           error
	switchToList  func() tea.Model
	currentConfig config.Config
}

type (
	errMsg error
)

// ValidateInputs checks if the provided MAC address and IP address are valid
func validateInputs(macAddress, ipAddress string) bool {
	return macAddressRegex.MatchString(macAddress) && ipAddressRegex.MatchString(ipAddress)
}

// InitialModel returns the initial model for the Device component
func InitialModel(switchToList func() tea.Model) Model {

	m := Model{
		err:           nil,
		switchToList:  switchToList,
		inputs:        make([]textinput.Model, 4), // Initialize the slice with length 4
		currentConfig: config.ReadConfig(),
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
			ti.Placeholder = "Enter the device name"
			ti.Focus()
			ti.PromptStyle = style.FocusedStyle
			ti.TextStyle = style.FocusedStyle
		// Description
		case 1:
			ti.Placeholder = "Enter a description for the device"
		// MAC address
		case 2:
			ti.Placeholder = "00:00:00:00:00:00"
		// IP address
		case 3:
			ti.Placeholder = "0.0.0.0"
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
		switch msg.Type {
		// Return to the list
		case tea.KeyEsc:
			return m.switchToList(), nil
		// Exit the program
		case tea.KeyCtrlC:
			return m, tea.Quit

		// Set focus to next input
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyDown, tea.KeyUp:
			// Check if the user pressed enter with the submit button focused
			if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				if m.focusIndex == len(m.inputs) {
					// Validate inputs
					if !validateInputs(m.inputs[2].Value(), m.inputs[3].Value()) {
						// Prevent the user from submitting the form
						return m, func() tea.Msg {
							return errMsg(fmt.Errorf("invalid MAC address or IP address"))
						}
					}

					// Append the device to the config
					updatedDevices := append(m.currentConfig.Devices, config.Device{
						DeviceName:  m.inputs[0].Value(),
						Description: m.inputs[1].Value(),
						MacAddress:  m.inputs[2].Value(),
						IPAddress:   m.inputs[3].Value(),
					})

					// Create a new config with the updated devices
					updatedConfig := config.Config{
						Devices: updatedDevices,
					}

					// Write the the new version of the config to the file
					config.WriteConfig(updatedConfig)

					// Return to the list and clear the screen
					return m.switchToList(), func() tea.Msg {
						return tea.ClearScreen()
					}
				}
			}

			// Cycle indexes
			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
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

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd

}

// updateInputs updates all the text inputs in the Device model.
func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View function for the Device model
func (m Model) View() string {
	// The header
	s := "New Device\n\n"

	// Add the text input fields
	for _, input := range m.inputs {
		s += input.View() + "\n"
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	s += fmt.Sprintf("\n\n%s\n\n", *button)

	s += style.HelpStyle.Render("\nPress esc to return to the list")

	return s

}
