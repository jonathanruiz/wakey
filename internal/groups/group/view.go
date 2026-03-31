package group

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
		key.WithKeys("q", "esc", "v"),
		key.WithHelp("q/esc", "back"),
	),
}

type DetailModel struct {
	group         config.Group
	deviceNames   []string
	previousModel tea.Model
	help          help.Model
	keys          detailKeyMap
}

func InitialDetailModel(previousModel tea.Model, selectedRow []string) DetailModel {
	cfg := config.ReadConfig()

	// Resolve device IDs to names
	deviceIDMap := make(map[string]string)
	for _, d := range cfg.Devices {
		deviceIDMap[d.ID] = d.DeviceName
	}

	// Find the group from config to get device IDs
	var grp config.Group
	for _, g := range cfg.Groups {
		if g.ID == selectedRow[0] {
			grp = g
			break
		}
	}

	var deviceNames []string
	for _, id := range grp.Devices {
		if name, ok := deviceIDMap[id]; ok {
			deviceNames = append(deviceNames, name)
		}
	}

	return DetailModel{
		group:         grp,
		deviceNames:   deviceNames,
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
	header := style.FocusedTab.Render("Groups > " + m.group.GroupName)
	s := lipgloss.PlaceHorizontal(style.TermWidth, lipgloss.Center, header) + "\n\n"

	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("98")).Width(16)

	devicesVal := "None"
	if len(m.deviceNames) > 0 {
		devicesVal = strings.Join(m.deviceNames, ", ")
	}

	s += fmt.Sprintf("%s %s\n", label.Render("Group Name"), m.group.GroupName)
	s += fmt.Sprintf("%s %s\n", label.Render("Devices"), devicesVal)

	s += "\n" + m.help.View(m.keys)

	return tea.NewView(s)
}
