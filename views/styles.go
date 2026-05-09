package views

import "github.com/charmbracelet/lipgloss"

var (
	PanelBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)

	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1)

	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	MutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	OkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	WarnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("31")).Bold(true)
	MarkedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("58")).Bold(true)
	InputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("25")).Bold(true)

	PlanColor1 = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	PlanColor2 = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	PlanColor3 = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	PlanColor4 = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	PlanColor5 = lipgloss.NewStyle().Foreground(lipgloss.Color("111")).Bold(true)
	PlanColor6 = lipgloss.NewStyle().Foreground(lipgloss.Color("190")).Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("24")).
			Padding(0, 1)
)
