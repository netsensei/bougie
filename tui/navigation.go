package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/tui/constants"
)

type Navigation struct {
	input textinput.Model
	mode  mode
}

func NewNavigation() Navigation {
	input := textinput.New()
	input.Prompt = "Bougie > "
	input.Placeholder = "go to..."
	input.CharLimit = 250
	input.Focus()

	m := Navigation{
		mode:  nav,
		input: input,
	}

	return m
}

func (m Navigation) Init() tea.Cmd {
	return nil
}

func (m Navigation) Update(msg tea.Msg) (Navigation, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case GopherQueryMsg:
		m.input.SetValue(msg.url)

	case ModeMsg:
		switch mode(msg) {
		case view:
			m.input.Blur()
		case nav:
			m.input.Focus()
		}

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.View) {
			cmds = append(cmds, SetBrowserModeCmd(view))
		}

		if key.Matches(msg, constants.Keymap.Enter) {
			value := m.input.Value()
			if value != "" {
				cmds = append(cmds, AddHistoryCmd(value))
				cmds = append(cmds, StartQueryCmd(value))
				return m, tea.Batch(cmds...)
			}
		}

		m.input, cmd = m.input.Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Navigation) View() string {
	return m.input.View()
}
