package popup

import (
	"fmt"
	"wakey/internal/status"
	"wakey/internal/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	yesFocusedButton = style.FocusedStyle.Render("[ (Y)es ]")                    // The focused button
	yesBlurredButton = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("(Y)es")) // The blurred button
	noFocusedButton  = style.FocusedStyle.Render("[ (N)o ]")                     // The focused button
	noBlurredButton  = fmt.Sprintf("[ %s ]", style.BlurredStyle.Render("(N)o"))  // The blurred button
)

type PopupMsg struct {
	message       string
	previousModel tea.Model
	handleFunc    HandleFunc
	help          help.Model
	table         table.Model
	focusIndex    int
	keyMap        keyMap
}

type HandleFunc func(selectedRow []string) (string, error)

func NewPopupMsg(message string, previousModel tea.Model, table table.Model, handleFunc HandleFunc) PopupMsg {
	return PopupMsg{
		message:       message,
		previousModel: previousModel,
		handleFunc:    handleFunc,
		table:         table,
		keyMap:        keys,
		help:          help.New(),
	}
}

func (m PopupMsg) Init() tea.Cmd { return nil }

func (m PopupMsg) handleYes(handleFunc HandleFunc) (tea.Model, tea.Cmd) {
	selected := m.table.SelectedRow()
	action, err := handleFunc(selected)
	if err != nil {
		status.Message = err
	} else {
		status.Message = fmt.Errorf(action)
	}

	return m.previousModel, func() tea.Msg {
		return tea.ClearScreen()
	}
}

func (m PopupMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Left):
			m.focusIndex = 0
		case key.Matches(msg, m.keyMap.Right):
			m.focusIndex = 1
		case key.Matches(msg, m.keyMap.Yes), key.Matches(msg, m.keyMap.Enter) && m.focusIndex == 0:
			return m.handleYes(m.handleFunc)
		case key.Matches(msg, m.keyMap.No), key.Matches(msg, m.keyMap.Enter) && m.focusIndex == 1:
			return m.previousModel, nil
		case key.Matches(msg, m.keyMap.Help):
			// Handle the "Help" key
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keyMap.Quit):
			return m.previousModel, nil
		}
	}

	var cmd tea.Cmd
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m PopupMsg) View() string {
	var buttons string

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(60).
		Height(5)

	if m.focusIndex == 0 {
		buttons = lipgloss.JoinHorizontal(lipgloss.Left, yesFocusedButton, noBlurredButton)
	} else {
		buttons = lipgloss.JoinHorizontal(lipgloss.Left, yesBlurredButton, noFocusedButton)
	}

	// Help text
	helpText := m.help.View(m.keyMap)

	return modalStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s", m.message, buttons, helpText))
}
