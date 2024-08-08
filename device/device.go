package device

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textInput    textinput.Model
	err          error
	switchToList func() tea.Model
}

type (
	errMsg error
)

func InitialModel(switchToList func() tea.Model) Model {
	ti := textinput.New()
	ti.Placeholder = "PC Name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
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

	// We handle errors just like any other message
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
	s += "\nPress q to quit.\n"

	return s
}
