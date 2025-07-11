package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Navigation struct {
	input textinput.Model
	Width int
	mode  mode
}

func NewNavigation() Navigation {
	input := textinput.New()
	input.Prompt = "> "
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

	case QueryMsg:
		m.input.Placeholder = msg.url

	case ModeMsg:
		m.mode = mode(msg)
		switch mode(msg) {
		case view:
			m.input.Blur()
		case nav:
			m.input.Focus()
		}

	case tea.KeyMsg:
		if m.mode == nav {
			m.input, cmd = m.input.Update(msg)
		}

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Navigation) View() string {
	return m.input.View()
}
