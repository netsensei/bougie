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
		m.status = ready
	case ModeMsg:
		m.mode = mode(msg)
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m Status) View() string {
	var status string
	var mode string
	//var action string

	if m.status == loading {
		status = fmt.Sprintf("%s loading %s...", m.spinner.View(), m.url)
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

	var statusBarStyle = lipgloss.NewStyle().
		Padding(1, 0).
		Bold(true)
	// Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
	// Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	var statusStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle)
		// Foreground(lipgloss.Color("#FFFDF5")).
		// Background(lipgloss.Color("#FF5F87")).
		// Padding(0, 1)

	var modeStyle = lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#6124DF")).
		Padding(0, 1).
		Align(lipgloss.Right)

	statusKey := statusStyle.Render(status)
	modeKey := modeStyle.Render(mode)
	midKey := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Width(m.Width - lipgloss.Width(statusKey) - lipgloss.Width(modeKey)).
		Render(" ")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		midKey,
		modeKey,
	)

	return statusBarStyle.Width(m.Width).Render(bar)
}
