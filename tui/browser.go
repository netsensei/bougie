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
	ready    bool
}

func initBrowser() (tea.Model, tea.Cmd) {
	status := NewStatus()
	navigation := NewNavigation()
	canvas := NewCanvas()

	m := Browser{
		// mode:       nav,
		status:     status,
		navigation: navigation,
		canvas:     canvas,
	}

	return m, func() tea.Msg { return nil }
}

func (m Browser) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.status.Init())

	cmds = append(cmds, SetBrowserModeCmd(view))
	cmds = append(cmds, AddHistoryCmd("gopher://floodgap.com"))
	cmds = append(cmds, StartQueryCmd("gopher://floodgap.com"))

	return tea.Batch(cmds...)
}

func (m Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}

		if m.ready {
			constants.WindowHeight = msg.Height - lipgloss.Height(m.navigation.View()) - lipgloss.Height(m.status.View())
			constants.WindowWidth = msg.Width
		}

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

var (
	InputStyle = lipgloss.NewStyle()
)

func (m Browser) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	return InputStyle.Render(m.navigation.View() + "\n" + m.canvas.View() + "\n" + m.status.View())
}
