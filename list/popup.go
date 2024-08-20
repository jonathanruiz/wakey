package list

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PopupMsg is not being used but will be retained for future use

type PopupMsg struct {
	message       string
	previousModel tea.Model
	timer         timer.Model
}

type timeoutMsg struct{}

func NewPopupMsg(message string, previousModel tea.Model) PopupMsg {
	t := timer.NewWithInterval(3*time.Second, time.Second)
	return PopupMsg{message: message, previousModel: previousModel, timer: t}
}

func (m PopupMsg) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return timeoutMsg{}
	}))
}

func (m PopupMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "esc" {
			return m.previousModel, nil
		}
	case timeoutMsg:
		return m.previousModel, nil
	}

	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m PopupMsg) View() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(50).
		Height(5)

	timeLeft := m.timer.View()
	return modalStyle.Render(fmt.Sprintf("%s\n\nTime left: %s", m.message, timeLeft))
}
