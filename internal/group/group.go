package group

import (
	"wakey/internal/config"
	"wakey/internal/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the Group component
type Model struct {
	groups []config.Group // list of groups
	keys   keyMap
	help   help.Model
	table  table.Model
}

// Init function for the Device model
func (m Model) Init() tea.Cmd { return nil }

// InitialModel function for the Group model
func InitialModel() tea.Model {
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
		devicesInGroup := ""
		for j, device := range group.Devices {
			devicesInGroup += device
			if j < len(group.Devices)-1 {
				devicesInGroup += ", "
			}
		}

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
		keys:  keys,
		help:  help.New(),
		table: t,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}
