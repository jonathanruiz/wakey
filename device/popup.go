package device

import (
	"fmt"
	"wakey/config"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PopupMsg struct {
	message       string
	previousModel tea.Model
	table         table.Model
}

func NewPopupMsg(message string, previousModel tea.Model, table table.Model) PopupMsg {
	return PopupMsg{message: message, previousModel: previousModel, table: table}
}

func (m PopupMsg) Init() tea.Cmd { return nil }

func (m PopupMsg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "y" {
			// Get the selected device
			selected := m.table.SelectedRow()

			// Get the current config
			currentConfig := config.ReadConfig()

			// Delete the selected device from the config
			for i, device := range currentConfig.Devices {
				if device.DeviceName == selected[0] {
					currentConfig.Devices = append(currentConfig.Devices[:i], currentConfig.Devices[i+1:]...)
					break
				}
			}

			// Write the new config to the file
			config.WriteConfig(currentConfig)

			return m.previousModel, nil
		}
		if msg.String() == "n" || msg.String() == "esc" {
			return m.previousModel, nil
		}
	}

	var cmd tea.Cmd
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

	return modalStyle.Render(fmt.Sprintf("%s\n", m.message))
}
