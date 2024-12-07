package main

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	foo      foo
	quitting bool
	err      error
}

func initialModel() model {
	f := NewFoo()
	return model{foo: f}
}

func (m model) Init() tea.Cmd {
	return m.foo.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.foo, cmd = m.foo.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := m.foo.View()

	if m.quitting {
		return str + "\n"
	}
	return str
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type foo struct {
	spinner spinner.Model
}

func NewFoo() foo {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return foo{
		spinner: s,
	}
}

func (m foo) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m foo) Update(msg tea.Msg) (foo, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m foo) View() string {
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	return str
}
