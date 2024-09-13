package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // default width if there's an error
	}
	return width
}

var TermWidth = GetTerminalWidth()

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
	countStyle          lipgloss.Style
	errStyle            lipgloss.Style
	statusStyle         lipgloss.Style
	statusMessageStyle  lipgloss.Style
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
	CountStyle          = CursorModeHelpStyle
	ErrStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("178")).Render // The error style
	StatusStyle         = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("241")).
				Width(TermWidth - 2)
	StatusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229")) // The status message style
	FocusedTab         = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("98")).Padding(0, 1)
	BlurredTab         = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("240")).Padding(0, 1)
)
