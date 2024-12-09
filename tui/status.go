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

	if m.status == loading {
		status = fmt.Sprintf("%s loading %s", m.spinner.View(), m.url)
	}

	if m.status == ready {
		status = "ready"
	}

	if m.mode == nav {
		mode = "INPUT"
	}
	if m.mode == view {
		mode = "READING"
	}

	return fmt.Sprintf("%s | %s", status, mode)
}
