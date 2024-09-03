package group

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the Group component
type Model struct {
	groups []string // list of groups
	keys   keyMap
	help   help.Model
	table  table.Model
}

// Init function for the Device model
func (m Model) Init() tea.Cmd { return nil }

// InitialModel function for the Group model
func InitialModel(previousModel tea.Model) tea.Model {
	// Define table columns
	columns := []table.Column{
		{Title: "Group", Width: 20},
	}

	// Define table rows
	rows := make([]table.Row, 0)

	// Create the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return Model{
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
