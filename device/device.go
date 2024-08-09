package device

import (
	"fmt"
	"wakey/config"
	"wakey/style"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedButton = style.FocusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Submit"))
)

// Model is the model for the Device component
type Model struct {
	focusIndex    int
	inputs        []textinput.Model
	cursorMode    cursor.Mode
	err           error
	switchToList  func() tea.Model
	addChoice     func(string)
	currentConfig config.Config
}

type (
	errMsg error
)

// InitialModel returns the initial model for the Device component
func InitialModel(switchToList func() tea.Model) Model {

	m := Model{
		err:           nil,
		switchToList:  switchToList,
		inputs:        make([]textinput.Model, 4), // Initialize the slice with length 4
		currentConfig: config.ReadConfig(),
	}

	var ti textinput.Model
	for i := range m.inputs {
		ti = textinput.New()
		ti.Cursor.Style = style.FocusedStyle
		ti.CharLimit = 64

		switch i {
		case 0:
			ti.Placeholder = "Enter the device name"
			ti.Focus()
			ti.PromptStyle = style.FocusedStyle
			ti.TextStyle = style.FocusedStyle
		case 1:
			ti.Placeholder = "Enter a description for the device"
		case 2:
			ti.Placeholder = "Enter the MAC address"
		case 3:
			ti.Placeholder = "Enter the IP address"
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

		case tea.KeyEsc:
			return m.switchToList(), nil
		case tea.KeyCtrlC:
			return m, tea.Quit

		// Change cursor mode
		case tea.KeyCtrlR:
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyDown, tea.KeyUp:

			// Check if the user pressed enter with the submit button focused
			if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				if m.focusIndex == len(m.inputs) {
					// Append the device to the config
					updatedDevices := append(m.currentConfig.Devices, m.inputs[0].Value())

					// Create a new config with the updated devices
					updatedConfig := config.Config{
						Devices: updatedDevices,
					}

					// Write the the new version of the config to the file
					config.WriteConfig(updatedConfig)

					// Return to the list
					return m.switchToList(), nil
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

	s += style.HelpStyle.Render("cursor mode is ")
	s += style.CursorModeHelpStyle.Render(m.cursorMode.String())
	s += style.HelpStyle.Render(" (ctrl+r to change style)")
	s += style.HelpStyle.Render("\nPress esc to return to the list")

	return s

}
