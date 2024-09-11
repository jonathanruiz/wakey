package internal

import (
	"wakey/internal/devices"
	"wakey/internal/groups"

	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	DevicesView View = iota
	GroupsView
)

type Model struct {
	CurrentView View
	Devices     devices.Model
	Groups      groups.Model
}

func InitialModel() Model {
	return Model{
		CurrentView: DevicesView,
		Devices:     devices.InitialModel().(devices.Model),
		Groups:      groups.InitialModel().(groups.Model),
	}
}

func (m Model) Init() tea.Cmd {
	// Perform any initial setup here, if needed
	return nil
}

func (m *Model) SwitchView(view View) {
	m.CurrentView = view
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.CurrentView {
	case DevicesView:
		var cmd tea.Cmd
		updatedModel, cmd := m.Devices.Update(msg)
		m.Devices = updatedModel.(devices.Model)
		return m, cmd
	case GroupsView:
		var cmd tea.Cmd
		updatedModel, cmd := m.Groups.Update(msg)
		m.Groups = updatedModel.(groups.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.CurrentView {
	case DevicesView:
		return m.Devices.View()
	case GroupsView:
		return m.Groups.View()
	default:
		return ""
	}
}
