package tui

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/tui/constants"
)

type Canvas struct {
	viewport   viewport.Model
	ready      bool
	mode       mode
	doc        string
	content    string
	scheme     string
	currentUrl string
	vpOffset   int
	links      Links
}

func NewCanvas() Canvas {
	return Canvas{}
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
			c.viewport = viewport.New(viewport.WithWidth(constants.WindowWidth), viewport.WithHeight(constants.WindowHeight))
			c.viewport.KeyMap = CanvasKeyMap()
			c.ready = true
			c.mode = view
		} else {
			c.viewport.SetWidth(constants.WindowWidth)
			c.viewport.SetHeight(constants.WindowHeight)
		}

	case ModeMsg:
		c.mode = mode(msg)

	case ReadyMsg:
		c.doc = msg.doc
		c.currentUrl = msg.currentUrl
		c.links = NewLinks(msg.links)
		c.content = msg.content
		c.scheme = msg.scheme

		c.viewport.SetContent(string(msg.content))
		cmds = append(cmds, SetBrowserModeCmd(view))

		return c, tea.Batch(cmds...)

	case CancelSearchMsg:
		c.viewport.SetContent(c.content)
		c.viewport.SetYOffset(c.vpOffset)

		cmds = append(cmds, SetBrowserModeCmd(view))
		return c, tea.Batch(cmds...)

	case ViewSourceMsg:
		c.viewport.SetContent(c.doc)
		c.viewport.SetYOffset(0)

	case RedrawMsg:
		c.viewport.SetContent(msg.content)
		offset := msg.active - (c.viewport.Height() / 2)
		c.viewport.SetYOffset(offset)

	case tea.KeyPressMsg:
		if key.Matches(msg, config.Keymap.Nav) {
			cmds = append(cmds, SetBrowserModeCmd(nav))
		}

		if key.Matches(msg, config.Keymap.View) {
			cmds = append(cmds, RedrawCmd(c.scheme, c.currentUrl, c.doc, c.links.ActiveLineNumber()))
			cmds = append(cmds, SetBrowserModeCmd(view))
		}

		if key.Matches(msg, config.Keymap.Source) {
			cmds = append(cmds, ViewSourceCmd())
			cmds = append(cmds, SetBrowserModeCmd(source))
		}

		if key.Matches(msg, config.Keymap.ItemForward) {
			if c.links.Forward() {
				cmds = append(cmds, RedrawCmd(c.scheme, c.currentUrl, c.doc, c.links.ActiveLineNumber()))
				return c, tea.Batch(cmds...)
			}
		}

		if key.Matches(msg, config.Keymap.ItemBackward) {
			if c.links.Backward() {
				cmds = append(cmds, RedrawCmd(c.scheme, c.currentUrl, c.doc, c.links.ActiveLineNumber()))
				return c, tea.Batch(cmds...)
			}
		}

		if key.Matches(msg, config.Keymap.Enter) {
			url := c.links.ActiveURL()
			if url == "" {
				return c, nil
			}
			cmds = append(cmds, AddHistoryCmd(url))
			cmds = append(cmds, StartQueryCmd(url))
			return c, tea.Batch(cmds...)
		}

		c.viewport, cmd = c.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c Canvas) View() tea.View {
	return tea.NewView(c.viewport.View())
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
