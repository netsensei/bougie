package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/tui/constants"
)

type Canvas struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func NewCanvas() Canvas {
	c := Canvas{
		content: "Bougie, a tiny sparking Gopher browser",
	}

	return c
}

func (m Canvas) Init() tea.Cmd {
	return nil
}

func (c Canvas) Update(msg tea.Msg) (Canvas, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !c.ready {
			c.viewport = viewport.New(constants.WindowWidth, constants.WindowHeight)
			c.viewport.SetContent(c.content)
			c.ready = true
		} else {
			c.viewport.Width = constants.WindowWidth
			c.viewport.Height = constants.WindowHeight
		}

	case ReadyMsg:
		c.content = string(msg)
		c.viewport.SetContent(c.content)
	}

	c.viewport, cmd = c.viewport.Update(msg)

	return c, cmd
}

func (c Canvas) View() string {
	return c.viewport.View()
}