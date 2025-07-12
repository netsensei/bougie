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
	Quit    key.Binding
	Nav     key.Binding
	View    key.Binding
	Enter   key.Binding
	Tab     key.Binding
	BackTab key.Binding
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
	Nav: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "Toggle nav mode"),
	),
	View: key.NewBinding(
		key.WithKeys("ctrl+v"),
		key.WithHelp("ctrl+v", "Toggle view mode"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Query"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next item"),
	),
	BackTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Previous item"),
	),
}
