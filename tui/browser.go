package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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

type Model struct {
	input    textinput.Model
	status   Status
	viewport viewport.Model
	mode     mode
	quitting bool
}

func initBrowser() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "go to..."
	input.CharLimit = 250
	input.Focus()

	status := NewStatus()

	m := Model{
		mode:   view,
		input:  input,
		status: status,
	}

	return m, func() tea.Msg { return nil }
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.status.Init())
	cmds = append(cmds, load("floodgap.com"))

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		height := msg.Height - lipgloss.Height(m.input.View()) - lipgloss.Height(m.status.View()) - 2
		m.viewport = viewport.New(msg.Width, height)
		m.input.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			m.quitting = true
			return m, tea.Quit
		}

		if m.mode == nav {
			if key.Matches(msg, constants.Keymap.Back) {
				m.mode = view
				m.input.Blur()
			}
			m.input, cmd = m.input.Update(msg)
		} else {
			if key.Matches(msg, constants.Keymap.Nav) {
				m.mode = nav
				m.input.Focus()
			}
			m.viewport, cmd = m.viewport.Update(msg)
		}
		cmds = append(cmds, cmd)

	case LoadingMsg:
		m.status, cmd = m.status.Update(msg)
		cmds = append(cmds, cmd)
		cmds = append(cmds, GetContent)
	case ReadyMsg:
		m.status, cmd = m.status.Update(msg)
		m.viewport.SetContent(string(msg))
		cmds = append(cmds, cmd)
	default:
		m.status, cmd = m.status.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	return constants.InputStyle.Render(m.input.View() + "\n" + m.viewport.View() + "\n" + m.status.View())
}
