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
	mode    status
	url     string
}

func NewStatus() Status {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Status{
		mode:    ready,
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
	case LoadingMsg:
		m.mode = loading
		m.url = msg.url
	case ReadyMsg:
		m.mode = ready
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m Status) View() string {
	if m.mode == loading {
		return fmt.Sprintf("%s loading %s", m.spinner.View(), m.url)
	}

	if m.mode == ready {
		return fmt.Sprintf("ready.")
	}

	return ""
}
