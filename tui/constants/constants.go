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
var InputStyle = lipgloss.NewStyle()

type keymap struct {
	Quit  key.Binding
	Mode  key.Binding
	Enter key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+q"),
		key.WithHelp("ctrl+c/ctrl+q", "Quit"),
	),
	// View: key.NewBinding(
	// 	key.WithKeys("esc"),
	// 	key.WithHelp("esc", "back"),
	// ),
	Mode: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "Toggle browser mode"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Query"),
	),
}
