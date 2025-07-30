package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

/* CONSTANTS */

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window
	WindowWidth  int
	WindowHeight int
)

type keymap struct {
	Quit         key.Binding
	Nav          key.Binding
	View         key.Binding
	Home         key.Binding
	Enter        key.Binding
	Tab          key.Binding
	BackTab      key.Binding
	PageForward  key.Binding
	PageBackward key.Binding
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
	Home: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "Back to the startpage"),
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
	PageForward: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Next page"),
	),
	PageBackward: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Previous page"),
	),
}
