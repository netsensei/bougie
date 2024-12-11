package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/tui/constants"
)

// TODO
//  * Build a status bar
//  * Show a handy help bar
//  * Wire up URL parsing CMD
//  * Wire up Query CMD with error msg & response msg
//  * Wire up Parsing cmd based on type of gopher content
//    * Pass gopher request document type via tea.Msg property

type mode int

const (
	nav mode = iota
	view
)

type status int

const (
	ready status = iota
	loading
)

type Browser struct {
	input  textinput.Model
	status Status
	// viewport viewport.Model
	canvas   Canvas
	mode     mode
	quitting bool
	// ready    bool
}

func initBrowser() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "go to..."
	input.CharLimit = 250
	input.Focus()

	status := NewStatus()

	m := Browser{
		mode:   nav,
		input:  input,
		status: status,
	}

	return m, func() tea.Msg { return nil }
}

func (m Browser) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.status.Init())
	cmds = append(cmds, queryCmd("floodgap.com"))

	return tea.Batch(cmds...)
}

func (m Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowHeight = msg.Height - lipgloss.Height(m.input.View()) - lipgloss.Height(m.status.View()) - 2
		constants.WindowWidth = msg.Width

		m.canvas = NewCanvas()
		m.input.Width = msg.Width
		m.status.Width = msg.Width

	case QueryMsg:
		cmds = append(cmds, GetContent("foo"))

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			m.quitting = true
			return m, tea.Quit
		}

		if key.Matches(msg, constants.Keymap.Mode) {
			if m.mode == nav {
				m.mode = view
				m.input.Blur()
			} else {
				m.mode = nav
				m.input.Focus()
			}
		}
		cmds = append(cmds, setBrowserMode(m.mode))
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	m.status, cmd = m.status.Update(msg)
	cmds = append(cmds, cmd)

	m.canvas, cmd = m.canvas.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Browser) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	return constants.InputStyle.Render(m.input.View() + "\n" + m.canvas.View() + "\n" + m.status.View())
}
