package list

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PopupMsg struct {
	message       string
	previousModel tea.Model
}

type timeoutMsg struct{}

func NewPopupMsg(message string, previousModel tea.Model) PopupMsg {
	return PopupMsg{message: message, previousModel: previousModel}
}

func (m PopupMsg) Init() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return timeoutMsg{}
	})
}

func (m PopupMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "esc" {
			return m.previousModel, nil
		}
	case timeoutMsg:
		return m.previousModel, nil
	}
	return m, nil
}

func (m PopupMsg) View() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(50).
		Height(5)

	return modalStyle.Render(m.message)
}
