package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/history"
	"github.com/netsensei/bougie/tui/constants"
)

type Canvas struct {
	viewport viewport.Model
	ready    bool
	mode     mode
	doc      string
	links    []map[int]string
	active   int
	history  *history.History
}

func NewCanvas() Canvas {
	c := Canvas{
		mode: nav,
		history: &history.History{
			Position: 0,
			Length:   0,
		},
		//	content: "Bougie, a tiny sparking Gopher browser",
	}

	return c
}

func (m Canvas) Init() tea.Cmd {
	return nil
}

func (c Canvas) Update(msg tea.Msg) (Canvas, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !c.ready {
			c.viewport = viewport.New(constants.WindowWidth, constants.WindowHeight)
			// c.viewport.SetContent(c.content)
			c.ready = true
		} else {
			c.viewport.Width = constants.WindowWidth
			c.viewport.Height = constants.WindowHeight
			c.viewport, cmd = c.viewport.Update(msg)
		}

	case ReadyMsg:
		if msg.err == nil {
			c.doc = msg.doc
			c.links = msg.links
			c.active = 0

			if len(msg.links) > 0 {
				keys := []int{}
				for k := range msg.links[0] {
					keys = append(keys, k)
				}

				offset := keys[0] - (c.viewport.Height / 2)
				c.viewport.SetYOffset(offset)
			} else {
				c.active = -1 // No links available
			}

			c.viewport.SetContent(string(msg.content))
			cmds = append(cmds, SetBrowserModeCmd(view))
			return c, tea.Batch(cmds...)
		}

	case RedrawMsg:
		c.viewport.SetContent(msg.content)
		offset := msg.position - (c.viewport.Height / 2)
		c.viewport.SetYOffset(offset)

	case AddHistoryMsg:
		c.history.Add(msg.url)

	case ModeMsg:
		c.mode = mode(msg)

	case tea.KeyMsg:
		if c.mode == view {
			if key.Matches(msg, constants.Keymap.Tab) {
				if c.active < len(c.links)-1 {
					c.active++

					keys := []int{}
					for k := range c.links[c.active] {
						keys = append(keys, k)
					}

					cmds = append(cmds, RedrawCmd(c.doc, keys[0]))
					return c, tea.Batch(cmds...)
				}
			}

			if key.Matches(msg, constants.Keymap.BackTab) {
				if c.active > 0 {
					c.active--

					keys := []int{}
					for k := range c.links[c.active] {
						keys = append(keys, k)
					}

					cmds = append(cmds, RedrawCmd(c.doc, keys[0]))
					return c, tea.Batch(cmds...)
				}
			}

			if key.Matches(msg, constants.Keymap.Enter) {
				if len(c.links[c.active]) == 0 {
					return c, nil // No links to follow
				}

				keys := []int{}
				for k := range c.links[c.active] {
					keys = append(keys, k)
				}

				cmds = append(cmds, AddHistoryCmd(c.links[c.active][keys[0]]))
				cmds = append(cmds, StartQueryCmd(c.links[c.active][keys[0]]))
				return c, tea.Batch(cmds...)
			}

			if key.Matches(msg, constants.Keymap.PageBackward) {
				if c.history.Length > 0 {
					c.history.Backward()
					url := c.history.Current()
					if url != "" {
						cmds = append(cmds, StartQueryCmd(url))
						return c, tea.Batch(cmds...)
					}

				}
			}

			if key.Matches(msg, constants.Keymap.PageForward) {
				if c.history.Length > 0 {
					c.history.Forward()
					url := c.history.Current()
					if url != "" {
						cmds = append(cmds, StartQueryCmd(url))
						return c, tea.Batch(cmds...)
					}

				}
			}

			c.viewport, cmd = c.viewport.Update(msg)
		}
	}

	return c, cmd
}

func (c Canvas) View() string {
	return c.viewport.View()
}
