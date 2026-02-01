package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/spf13/viper"
)

type keymap struct {
	Quit         key.Binding
	Nav          key.Binding
	View         key.Binding
	Source       key.Binding
	Home         key.Binding
	Reload       key.Binding
	Enter        key.Binding
	ItemForward  key.Binding
	ItemBackward key.Binding
	CmpntForward key.Binding
	PageForward  key.Binding
	PageBackward key.Binding
	HalfPageDown key.Binding
	HalfPageUp   key.Binding
	PageDown     key.Binding
	PageUp       key.Binding
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
}

var Keymap keymap

func keysInit() {
	var bindings = map[string][]string{
		"Quit":         viper.GetStringSlice("keybindings.quit"),
		"Nav":          viper.GetStringSlice("keybindings.nav"),
		"View":         viper.GetStringSlice("keybindings.view"),
		"Source":       viper.GetStringSlice("keybindings.source"),
		"Home":         viper.GetStringSlice("keybindings.home"),
		"Reload":       viper.GetStringSlice("keybindings.reload"),
		"Enter":        viper.GetStringSlice("keybindings.enter"),
		"ItemForward":  viper.GetStringSlice("keybindings.item_forward"),
		"ItemBackward": viper.GetStringSlice("keybindings.item_backward"),
		"CmpntForward": viper.GetStringSlice("keybindings.component_forward"),
		"PageForward":  viper.GetStringSlice("keybindings.page_forward"),
		"PageBackward": viper.GetStringSlice("keybindings.page_backward"),
		"PageDown":     viper.GetStringSlice("keybindings.page_down"),
		"PageUp":       viper.GetStringSlice("keybindings.page_up"),
		"Up":           viper.GetStringSlice("keybindings.up"),
		"Down":         viper.GetStringSlice("keybindings.down"),
		"Left":         viper.GetStringSlice("keybindings.left"),
		"Right":        viper.GetStringSlice("keybindings.right"),
	}

	// Keymap reusable key mappings shared across models
	Keymap = keymap{
		Quit: key.NewBinding(
			key.WithKeys(bindings["Quit"]...),
			key.WithHelp(strings.Join(bindings["Quit"], ", "), "Quit"),
		),
		Nav: key.NewBinding(
			key.WithKeys(bindings["Nav"]...),
			key.WithHelp(strings.Join(bindings["Nav"], ", "), "Toggle nav mode"),
		),
		View: key.NewBinding(
			key.WithKeys(bindings["View"]...),
			key.WithHelp(strings.Join(bindings["View"], ", "), "Toggle view mode"),
		),
		Source: key.NewBinding(
			key.WithKeys(bindings["Source"]...),
			key.WithHelp(strings.Join(bindings["Source"], ", "), "Toggle source mode"),
		),
		Home: key.NewBinding(
			key.WithKeys(bindings["Home"]...),
			key.WithHelp(strings.Join(bindings["Home"], ", "), "Go back home"),
		),
		Reload: key.NewBinding(
			key.WithKeys(bindings["Reload"]...),
			key.WithHelp(strings.Join(bindings["Reload"], ", "), "Reload the current resource"),
		),
		Enter: key.NewBinding(
			key.WithKeys(bindings["Enter"]...),
			key.WithHelp(strings.Join(bindings["Enter"], ", "), "Query"),
		),
		ItemForward: key.NewBinding(
			key.WithKeys(bindings["ItemForward"]...),
			key.WithHelp(strings.Join(bindings["ItemForward"], ", "), "Next item"),
		),
		ItemBackward: key.NewBinding(
			key.WithKeys(bindings["ItemBackward"]...),
			key.WithHelp(strings.Join(bindings["ItemBackward"], ", "), "Previous item"),
		),
		PageForward: key.NewBinding(
			key.WithKeys(bindings["PageForward"]...),
			key.WithHelp(strings.Join(bindings["PageForward"], ", "), "Next page"),
		),
		PageBackward: key.NewBinding(
			key.WithKeys(bindings["PageBackward"]...),
			key.WithHelp(strings.Join(bindings["PageBackward"], ", "), "Previous page"),
		),
		CmpntForward: key.NewBinding(
			key.WithKeys(bindings["CmpntForward"]...),
			key.WithHelp(strings.Join(bindings["CmpntForward"], ", "), "Select next form element"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys(bindings["HalfPageDown"]...),
			key.WithHelp(strings.Join(bindings["HalfPageDown"], ", "), "Scroll half page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys(bindings["HalfPageUp"]...),
			key.WithHelp(strings.Join(bindings["HalfPageUp"], ", "), "Scroll half page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys(bindings["PageDown"]...),
			key.WithHelp(strings.Join(bindings["PageDown"], ", "), "Scroll one page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys(bindings["PageUp"]...),
			key.WithHelp(strings.Join(bindings["PageUp"], ", "), "Scroll one page up"),
		),
		Up: key.NewBinding(
			key.WithKeys(bindings["Up"]...),
			key.WithHelp(strings.Join(bindings["Up"], ", "), "Scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys(bindings["Down"]...),
			key.WithHelp(strings.Join(bindings["Down"], ", "), "Scroll down"),
		),
		Left: key.NewBinding(
			key.WithKeys(bindings["Left"]...),
			key.WithHelp(strings.Join(bindings["Left"], ", "), "Scroll left"),
		),
		Right: key.NewBinding(
			key.WithKeys(bindings["Right"]...),
			key.WithHelp(strings.Join(bindings["Right"], ", "), "Scroll right"),
		),
	}

}
