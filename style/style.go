package style

import "github.com/charmbracelet/lipgloss"

// Define the style struct
type Style struct {
	focusedStyle        lipgloss.Style
	blurredStyle        lipgloss.Style
	cursorStyle         lipgloss.Style
	noStyle             lipgloss.Style
	helpStyle           lipgloss.Style
	cursorModeHelpStyle lipgloss.Style
	titleStyle          lipgloss.Style
	italicHeaderStyle   lipgloss.Style
	deviceCountStyle    lipgloss.Style
}

// Define the style struct
type TableStyles struct {
	Header   lipgloss.Style
	Selected lipgloss.Style
}

// DefaultTableStyles returns a set of default style definitions for the table.
func DefaultTableStyles() TableStyles {
	return TableStyles{
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Foreground(lipgloss.Color("98")),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("98")).
			Bold(false),
	}
}

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("98"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	TitleStyle          = lipgloss.NewStyle().Bold(true)
	ItalicHeaderStyle   = lipgloss.NewStyle().Italic(true)
	DeviceCountStyle    = CursorModeHelpStyle
)
