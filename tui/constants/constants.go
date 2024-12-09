package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* CONSTANTS */

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window
	WindowWidth  int
	WindowHeight int
)
var InputStyle = lipgloss.NewStyle().Margin(0, 2)

type keymap struct {
	Quit  key.Binding
	View  key.Binding
	Nav   key.Binding
	Enter key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+q"),
		key.WithHelp("ctrl+c/ctrl+q", "quit"),
	),
	View: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Nav: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "Navigate"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "browse"),
	),
}
