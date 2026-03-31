package device

import (
	"fmt"
	"strings"
	"wakey/internal/common/style"
	"wakey/internal/config"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type detailKeyMap struct {
	Quit key.Binding
}

func (k detailKeyMap) ShortHelp() []key.Binding { return []key.Binding{k.Quit} }
func (k detailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

var detailKeys = detailKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "i"),
		key.WithHelp("q/esc", "back"),
	),
}

type DetailModel struct {
	device        config.Device
	groups        []string
	previousModel tea.Model
	help          help.Model
	keys          detailKeyMap
}

func InitialDetailModel(previousModel tea.Model, selectedRow []string) DetailModel {
	cfg := config.ReadConfig()

	device := config.Device{
		ID:          selectedRow[0],
		DeviceName:  selectedRow[1],
		Description: selectedRow[2],
		MacAddress:  selectedRow[3],
		IPAddress:   selectedRow[4],
		State:       selectedRow[5],
	}

	var groupNames []string
	for _, g := range cfg.Groups {
		for _, id := range g.Devices {
			if id == device.ID {
				groupNames = append(groupNames, g.GroupName)
				break
			}
		}
	}

	return DetailModel{
		device:        device,
		groups:        groupNames,
		previousModel: previousModel,
		help:          help.New(),
		keys:          detailKeys,
	}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m.previousModel, nil
		}
	}
	return m, nil
}

func (m DetailModel) View() tea.View {
	header := style.FocusedTab.Render("Devices > " + m.device.DeviceName)
	s := lipgloss.PlaceHorizontal(style.TermWidth, lipgloss.Center, header) + "\n\n"

	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("98")).Width(16)

	s += fmt.Sprintf("%s %s\n", label.Render("Device Name"), m.device.DeviceName)
	s += fmt.Sprintf("%s %s\n", label.Render("Description"), m.device.Description)
	s += fmt.Sprintf("%s %s\n", label.Render("MAC Address"), m.device.MacAddress)
	s += fmt.Sprintf("%s %s\n", label.Render("IP Address"), m.device.IPAddress)
	s += fmt.Sprintf("%s %s\n", label.Render("State"), m.device.State)

	groupsVal := "None"
	if len(m.groups) > 0 {
		groupsVal = strings.Join(m.groups, ", ")
	}
	s += fmt.Sprintf("%s %s\n", label.Render("Groups"), groupsVal)

	s += "\n" + m.help.View(m.keys)

	return tea.NewView(s)
}
