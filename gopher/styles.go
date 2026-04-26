package gopher

import (
	"charm.land/lipgloss/v2"
	"github.com/netsensei/bougie/tui/constants"
)

var documentStyle = lipgloss.NewStyle()

//	Background(lipgloss.Color("#7D56F4"))

var textStyle = lipgloss.NewStyle().
	Inherit(documentStyle).
	Width(constants.WindowWidth).
	Foreground(lipgloss.Color("#FAFAFA"))

var typeStyle = lipgloss.NewStyle().
	Inherit(documentStyle).
	Width(6)

var linkStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))

var activeLinkStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Bold(true).
	Foreground(lipgloss.Color("#CC56F4"))
