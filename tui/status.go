package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var style = lipgloss.NewStyle().
// 	Bold(true)

type Status struct {
	spinner spinner.Model
	status  status
	mode    mode
	url     string
	Width   int
	err     error
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
	case QueryMsg:
		m.status = loading
		m.url = msg.url
	case ReadyMsg:
		if msg.err != nil {
			m.status = errored
			m.url = msg.url
			m.err = msg.err
		} else {
			m.status = ready
			m.err = nil
		}
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
		statusMsg = fmt.Sprintf("loading %s...", m.url)

	}

	if m.status == saving {
		statusMsg = fmt.Sprintf("saving %s...", m.url)
	}

	if m.status == errored {
		status = fmt.Sprintf("Error: %v", m.err)
	}

	if m.status == ready {
		status = "Ready."
	}

	if m.mode == nav {
		mode = "Input"
	}
	if m.mode == view {
		mode = "View"
	}

	barStyle := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#6124DF"})

	statusMsgStyle := lipgloss.NewStyle().
		Inherit(barStyle)

	if m.status == saving || m.status == loading {
		statusKey := statusMsgStyle.Render(statusMsg)
		status = lipgloss.JoinHorizontal(lipgloss.Top, m.spinner.View(), statusKey)
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
