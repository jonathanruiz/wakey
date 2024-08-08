package device

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	viewport     viewport.Model
	textInput    textinput.Model
	err          error
	switchToList func() tea.Model
}

type (
	errMsg error
)

func InitialModel(switchToList func() tea.Model) Model {
	vp := viewport.New(20, 10) // Adjust width and height as needed

	// Create a new text input field
	ti := textinput.New()
	ti.Placeholder = "Enter a new device name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		viewport:     vp,
		textInput:    ti,
		err:          nil,
		switchToList: switchToList,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			return m.switchToList(), nil
		case tea.KeyEsc:
			return m.switchToList(), nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	// The header
	s := "New Device\n\n"
	s += m.textInput.View() + "\n"
	s += "\nPress esc to cancel."

	s += "\n" + m.viewport.View()
	return s
}
