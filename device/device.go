package device

import (
	"fmt"
	"wakey/config"
	"wakey/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	switchToList  func() tea.Model
	currentConfig config.Config
}

// InitialModel returns the initial model for the Device component
func InitialModel(switchToList func() tea.Model) Model {

	m := Model{
		err:           make([]error, 4), // Initialize the slice with length 4
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

	// Ensure focus is correctly managed
	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex].Focus()
	}

	return tea.Batch(cmds...)
}

// View function for the Device model
func (m Model) View() string {
	// The header
	s := "New Device\n\n"

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

	s += style.HelpStyle.Render("\nPress esc to return to the list")

	return s

}
