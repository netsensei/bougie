package tui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorPrimary    = lipgloss.Color("#F25D94")
	ColorSecondary  = lipgloss.Color("#874BFD")
	ColorForeground = lipgloss.Color("#FFF7DB")
	ColorBackground = lipgloss.Color("#555555")
	ColorSubtle     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorBarBg      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#6124DF"}
	ColorSpinner    = lipgloss.Color("205")
	ColorButtonBg   = lipgloss.Color("#888B7E")
)

// Navigation styles
var (
	LogoStyle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorForeground).
			Align(lipgloss.Center).
			Width(12)

	NavBaseStyle = lipgloss.NewStyle().
			Foreground(ColorForeground).
			Background(ColorBackground).
			Padding(0, 1)
)

// Status bar styles
var (
	BarStyle = lipgloss.NewStyle().
			Background(ColorBarBg)

	StatusMsgStyle = lipgloss.NewStyle().
			Inherit(BarStyle)

	StatusKeyStyle = lipgloss.NewStyle().
			Inherit(BarStyle).
			Padding(0, 1)

	ModeStyle = lipgloss.NewStyle().
			Inherit(BarStyle).
			Width(10).
			Align(lipgloss.Center)
)

// Search dialog styles
var (
	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(ColorForeground).
			Background(ColorButtonBg).
			Padding(0, 3).
			MarginTop(1)

	ActiveButtonStyle = lipgloss.NewStyle().
				Foreground(ColorForeground).
				Background(ColorPrimary).
				Padding(0, 3).
				MarginTop(1).
				Underline(true)

	QuestionStyle = lipgloss.NewStyle().
			Width(75).
			Align(lipgloss.Center)
)

// Canvas styles
var (
	CanvasStyle = lipgloss.NewStyle().
		Padding(0, 1)
)
