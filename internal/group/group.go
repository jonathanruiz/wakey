package group

import (
	"wakey/internal/config"
	"wakey/internal/newGroup"
	"wakey/internal/status"
	"wakey/internal/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the Group component
type Model struct {
	groups        []config.Group // list of groups
	previousModel tea.Model
	keys          keyMap
	help          help.Model
	table         table.Model
}

// Init function for the Device model
func (m Model) Init() tea.Cmd { return nil }

// InitialModel function for the Group model
func InitialModel(previousModel tea.Model) tea.Model {
	// Get groups with updated state
	groups := config.GetUpdateState().Groups

	// Define table columns
	columns := []table.Column{
		{Title: "Group Name", Width: 20},
		{Title: "Devices", Width: 30},
	}

	// Define table rows
	rows := make([]table.Row, len(groups))
	for i, group := range groups {
		// Create a string of devices in the group
		devicesInGroup := ""
		for _, device := range group.Devices {
			devicesInGroup += device + ", "
		}

		// Remove the trailing comma and space
		devicesInGroup = devicesInGroup[:len(devicesInGroup)-2]

		rows[i] = table.Row{
			group.GroupName,
			devicesInGroup,
		}
	}

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set the custom key bindings
	t.KeyMap = table.KeyMap{
		LineUp:   keys.Up,
		LineDown: keys.Down,
	}

	// Get the default table styles
	s := style.DefaultTableStyles()

	// Set the styles
	t.SetStyles(table.Styles{
		Header:   s.Header,
		Selected: s.Selected,
	})

	return Model{
		// A list of devices to wake. This could be fetched from a database or config file
		groups: groups,
		// A map which indicates which devices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `devices` slice, above.
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
		table:         t,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Create):
			// Create a new group
			return newGroup.InitialModel(m), nil
		case key.Matches(msg, m.keys.Edit):
			// Edit the selected group
			break
		case key.Matches(msg, m.keys.Delete):
			// Delete the selected group
			break
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m.previousModel, tea.ClearScreen
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := "\n"
	s += m.table.View() + "\n"

	// Status message
	var statusMessage string
	if status.Message != nil {
		statusMessage = status.Message.Error()
	} else {
		statusMessage = "No status"
	}
	s += style.StatusStyle.Render("Status: "+style.StatusMessageStyle.Render(statusMessage)) + "\n"

	// Help text
	s += m.help.View(m.keys)
	return s
}
