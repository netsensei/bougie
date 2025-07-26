package tui

import (
	"net/url"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/tui/constants"
)

type searchDiagogActive int

const (
	input searchDiagogActive = iota
	cancel
	ok
)

type Search struct {
	input  textinput.Model
	mode   mode
	url    string
	Width  int
	Height int
	active searchDiagogActive
}

func NewSearch(Width int, Height int) Search {
	return Search{
		Width:  Width,
		Height: Height,
		active: input,
	}
}

func (m Search) Init() tea.Cmd {
	return nil
}

func (m Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ModeMsg:
		m.mode = mode(msg)

	case SearchMsg:
		input := textinput.New()
		input.Prompt = "Search > "
		input.Placeholder = "go to..."
		input.CharLimit = 250
		input.Focus()

		m.input = input
		m.url = msg.url

		cmds = append(cmds, SetBrowserModeCmd(search))
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if m.mode == search {
			if key.Matches(msg, constants.Keymap.Enter) {
				value := m.input.Value()

				if value != "" && m.active == ok {
					purl, _ := url.Parse(m.url)
					q := purl.Query()
					q.Set("q", value)
					purl.RawQuery = q.Encode()

					cmds = append(cmds, AddHistoryCmd(purl.String()))
					cmds = append(cmds, StartQueryCmd(purl.String()))
					return m, tea.Batch(cmds...)
				}

				if m.active == cancel {
					cmds = append(cmds, CancelSearchCmd())
					return m, tea.Batch(cmds...)
				}
			}

			if key.Matches(msg, constants.Keymap.Tab) {
				switch m.active {
				case input:
					m.input.Blur()
					m.active = ok
				case ok:
					m.input.Blur()
					m.active = cancel
				case cancel:
					m.input.Focus()
					m.active = input
				}
			}

			m.input, cmd = m.input.Update(msg)
		}

		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Search) View() string {
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	// Dialog.

	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1)

	activeButtonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		Padding(0, 3).
		MarginTop(1).
		Underline(true)

	var okButton, cancelButton string
	switch m.active {
	case input:
		okButton = buttonStyle.MarginRight(2).Render("Search")
		cancelButton = buttonStyle.Render("Cancel")
	case cancel:
		okButton = buttonStyle.MarginRight(2).Render("Search")
		cancelButton = activeButtonStyle.Render("Cancel")
	case ok:
		okButton = activeButtonStyle.MarginRight(2).Render("Search")
		cancelButton = buttonStyle.Render("Cancel")
	}

	question := lipgloss.NewStyle().Width(100).Align(lipgloss.Center).Render(m.input.View())
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("  "),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	return dialog
}
