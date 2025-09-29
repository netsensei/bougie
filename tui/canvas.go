package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/tui/constants"
)

type Canvas struct {
	viewport   viewport.Model
	ready      bool
	doc        string
	content    string
	scheme     string
	currentUrl string
	vpOffset   int
	links      []map[int]string
	active     int
}

func NewCanvas() Canvas {
	c := Canvas{
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
			c.viewport.KeyMap = CanvasKeyMap()
			c.ready = true
		} else {
			c.viewport.Width = constants.WindowWidth
			c.viewport.Height = constants.WindowHeight
		}

	case ReadyMsg:
		c.doc = msg.doc
		c.currentUrl = msg.currentUrl
		c.links = msg.links
		c.active = 0
		c.content = msg.content
		c.scheme = msg.scheme

		if len(msg.links) > 0 {
			keys := []int{}
			for k := range msg.links[0] {
				keys = append(keys, k)
			}

			offset := keys[0] - (c.viewport.Height / 2)

			c.vpOffset = offset
			c.viewport.SetYOffset(offset)
		} else {
			c.active = -1 // No links available
		}

		c.viewport.SetContent(string(msg.content))
		cmds = append(cmds, SetBrowserModeCmd(view))

		return c, tea.Batch(cmds...)

	case CancelSearchMsg:
		c.viewport.SetContent(c.content)
		c.viewport.SetYOffset(c.vpOffset)

		cmds = append(cmds, SetBrowserModeCmd(view))
		return c, tea.Batch(cmds...)

	case RedrawMsg:
		c.viewport.SetContent(msg.content)
		offset := msg.active - (c.viewport.Height / 2)
		c.viewport.SetYOffset(offset)

	case tea.KeyMsg:
		if key.Matches(msg, config.Keymap.Nav) {
			cmds = append(cmds, SetBrowserModeCmd(nav))
		}

		if key.Matches(msg, config.Keymap.ItemForward) {
			if c.active < len(c.links)-1 {
				c.active++

				keys := []int{}
				for k := range c.links[c.active] {
					keys = append(keys, k)
				}

				cmds = append(cmds, RedrawCmd(c.scheme, c.currentUrl, c.doc, keys[0]))
				return c, tea.Batch(cmds...)
			}
		}

		if key.Matches(msg, config.Keymap.ItemBackward) {
			if c.active > 0 {
				c.active--

				keys := []int{}
				for k := range c.links[c.active] {
					keys = append(keys, k)
				}

				cmds = append(cmds, RedrawCmd(c.scheme, c.currentUrl, c.doc, keys[0]))
				return c, tea.Batch(cmds...)
			}
		}

		if key.Matches(msg, config.Keymap.Enter) {
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

		c.viewport, cmd = c.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c Canvas) View() string {
	return c.viewport.View()
}

func CanvasKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		PageDown:     config.Keymap.PageDown,
		PageUp:       config.Keymap.PageUp,
		HalfPageUp:   config.Keymap.HalfPageUp,
		HalfPageDown: config.Keymap.HalfPageDown,
		Up:           config.Keymap.Up,
		Down:         config.Keymap.Down,
		Left:         config.Keymap.Left,
		Right:        config.Keymap.Right,
	}
}
