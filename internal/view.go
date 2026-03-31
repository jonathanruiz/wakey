package internal

import (
	"wakey/internal/common"
	"wakey/internal/devices"
	"wakey/internal/groups"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type View int

const (
	DevicesView View = iota
	GroupsView
)

type Model struct {
	CurrentView  View
	CurrentModel tea.Model
	Keys         common.KeyMap
}

func InitialModel() Model {
	return Model{
		CurrentView:  DevicesView,
		CurrentModel: devices.InitialModel(),
		Keys:         common.DefaultKeyMap(),
	}
}

func (m *Model) SwitchView(view View) {
	m.CurrentView = view
	switch view {
	case DevicesView:
		m.CurrentModel = devices.InitialModel()
	case GroupsView:
		m.CurrentModel = groups.InitialModel()
	}
}

func (m Model) Init() tea.Cmd {
	return m.CurrentModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.View):
			switch m.CurrentView {
			case DevicesView:
				m.SwitchView(GroupsView)
				return m, tea.ClearScreen
			case GroupsView:
				m.SwitchView(DevicesView)
				return m, tea.Batch(tea.ClearScreen, devices.UpdateStateCmd())
			}
		}
	}

	var cmd tea.Cmd
	m.CurrentModel, cmd = m.CurrentModel.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	v := m.CurrentModel.View()
	v.AltScreen = true
	return v
}
