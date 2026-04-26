package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Default is the built-in color palette. Future themes should satisfy the same
// shape so callers can swap them without changing rendering logic.
var Default = Palette{
	// Base colors
	Primary:    lipgloss.Color("#F25D94"),
	Secondary:  lipgloss.Color("#874BFD"),
	Foreground: lipgloss.Color("#FFF7DB"),
	Background: lipgloss.Color("#555555"),
	Subtle:     lipgloss.Color("#383838"),
	BarBg:      lipgloss.Color("#6124DF"),
	Spinner:    lipgloss.Color("205"),
	ButtonBg:   lipgloss.Color("#888B7E"),

	// Content colors
	ContentFg:    lipgloss.Color("#FAFAFA"),
	ContentMuted: lipgloss.Color("#AEAEAE"),
	Heading:      lipgloss.Color("#FFFFFF"),
	Link:         lipgloss.Color("#7D56F4"),
	LinkActive:   lipgloss.Color("#CC56F4"),
}

type Palette struct {
	// Base colors
	Primary    color.Color
	Secondary  color.Color
	Foreground color.Color
	Background color.Color
	Subtle     color.Color
	BarBg      color.Color
	Spinner    color.Color
	ButtonBg   color.Color

	// Content colors
	ContentFg    color.Color
	ContentMuted color.Color
	Heading      color.Color
	Link         color.Color
	LinkActive   color.Color
}
