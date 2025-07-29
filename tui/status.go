package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/tui/constants"
)

type Status struct {
	spinner  spinner.Model
	status   status
	mode     mode
	url      string
	resource string
	err      error
}

func NewStatus() Status {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Status{
		status:  ready,
		mode:    nav,
		spinner: s,
		url:     "",
	}

	return m
}

func (m Status) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Status) Update(msg tea.Msg) (Status, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case GopherDocumentQueryMsg:
		m.status = loading
		m.url = msg.url

	case GopherFileQueryCmd:
		m.status = saving
		m.url = msg.url

	case FileSavedMsg:
		m.status = saved
		m.resource = msg.resource

	case ErrorMsg:
		m.status = errored
		m.url = msg.url
		m.err = msg.err

	case ReadyMsg:
		m.status = ready
		m.err = nil

	case ModeMsg:
		m.mode = mode(msg)

	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m Status) View() string {
	var status string
	var statusMsg string
	var mode string

	if m.status == loading {
		statusMsg = fmt.Sprintf("Loading %s...", m.url)

	}

	if m.status == saving {
		statusMsg = fmt.Sprintf("Saving %s...", m.url)
	}

	if m.status == saved {
		statusMsg = fmt.Sprintf("Saved %s succesfully", m.resource)
	}

	if m.status == errored {
		status = fmt.Sprintf("Error: %v", m.err)
	}

	if m.status == ready {
		status = "Ready."
	}

	if m.mode == nav {
		mode = "Navigation"
	}
	if m.mode == view {
		mode = "View"
	}
	if m.mode == search {
		mode = "Search"
	}
	if m.mode == input {
		mode = "Input"
	}
	if m.mode == save {
		mode = "Saving"
	}

	barStyle := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#6124DF"})

	statusMsgStyle := lipgloss.NewStyle().
		Inherit(barStyle)

	if m.status == saving || m.status == loading {
		statusKey := statusMsgStyle.Render(statusMsg)
		status = lipgloss.JoinHorizontal(lipgloss.Top, m.spinner.View(), statusKey)
	}

	if m.status == saved {
		statusKey := statusMsgStyle.Render(statusMsg)
		status = lipgloss.JoinHorizontal(lipgloss.Top, statusKey)
	}

	statusStyle := lipgloss.NewStyle().
		Inherit(barStyle).
		Width(lipgloss.Width(status)+2).
		Padding(0, 1)

	modeStyle := lipgloss.NewStyle().
		Inherit(barStyle).
		Width(10).
		Align(lipgloss.Center)

	statusKey := statusStyle.Render(status)
	modeKey := modeStyle.Render(mode)
	midKey := lipgloss.NewStyle().
		Inherit(barStyle).
		Width(constants.WindowWidth - lipgloss.Width(statusKey) - lipgloss.Width(modeKey)).
		Render(" ")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, statusKey, midKey, modeKey)

	return bar
}
