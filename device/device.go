package device

import (
	"fmt"
	"wakey/style"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedButton = style.FocusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("Submit"))
)

type Model struct {
	viewport     viewport.Model
	focusIndex   int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	err          error
	switchToList func() tea.Model
	addChoice    func(string)
}

type (
	errMsg error
)

func InitialModel(switchToList func() tea.Model, addChoice func(string)) Model {
	vp := viewport.New(20, 10) // Adjust width and height as needed

	m := Model{
		viewport:     vp,
		err:          nil,
		switchToList: switchToList,
		addChoice:    addChoice,
		inputs:       make([]textinput.Model, 4), // Initialize the slice with length 4
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

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

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
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if m.focusIndex == len(m.inputs) {
					deviceName := m.inputs[0].Value()
					m.addChoice(deviceName)
					return m.switchToList(), nil
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

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

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

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

	s += m.viewport.View()

	return s

}
