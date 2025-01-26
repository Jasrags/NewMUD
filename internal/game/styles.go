package game

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1, 0, 1)

		// Table formatting
	singleColumnStyle = borderStyle.Width(80)
	dualColumnStyle   = borderStyle.Width(39)

	// Text styles
	headerStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	attrNameStyle      = lipgloss.NewStyle().Bold(true).Foreground(white)
	attrTextValueStyle = lipgloss.NewStyle().Bold(false).Foreground(white)
	attrValueStyle     = lipgloss.NewStyle().Bold(false).Foreground(cyan)
	attrPosModStyle    = lipgloss.NewStyle().Bold(true).Foreground(green)
	attrNegModStyle    = lipgloss.NewStyle().Bold(true).Foreground(red)

	// Colors
	black   = lipgloss.Color("0")
	red     = lipgloss.Color("1")
	green   = lipgloss.Color("2")
	yellow  = lipgloss.Color("3")
	blue    = lipgloss.Color("4")
	magenta = lipgloss.Color("5")
	cyan    = lipgloss.Color("6")
	white   = lipgloss.Color("7")
)
