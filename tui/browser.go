package tui

import (
	"github.com/charmbracelet/bubbles/key"
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
	saving
	errored
)

type Browser struct {
	status     Status
	canvas     Canvas
	navigation Navigation
	// mode       mode
	quitting bool
	// ready    bool
}

func initBrowser() (tea.Model, tea.Cmd) {
	status := NewStatus()
	navigation := NewNavigation()

	m := Browser{
		// mode:       nav,
		status:     status,
		navigation: navigation,
	}

	return m, func() tea.Msg { return nil }
}

func (m Browser) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.status.Init())

	return tea.Batch(cmds...)
}

func (m Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowHeight = msg.Height - lipgloss.Height(m.navigation.View()) - lipgloss.Height(m.status.View()) - 2
		constants.WindowWidth = msg.Width

		m.canvas = NewCanvas()

		cmds = append(cmds, AddHistoryCmd("gopher://floodgap.com"))
		cmds = append(cmds, StartQueryCmd("gopher://floodgap.com"))

		m.navigation.Width = msg.Width
		m.status.Width = msg.Width

	case QueryMsg:
		cmds = append(cmds, SendQueryCmd(msg.url))

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			m.quitting = true
			return m, tea.Quit
		}

		if key.Matches(msg, constants.Keymap.Nav) {
			cmds = append(cmds, SetBrowserModeCmd(nav))
		}

		if key.Matches(msg, constants.Keymap.View) {
			cmds = append(cmds, SetBrowserModeCmd(view))
		}
	}

	m.navigation, cmd = m.navigation.Update(msg)
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

	return constants.InputStyle.Render(m.navigation.View() + "\n" + m.canvas.View() + "\n" + m.status.View())
}
