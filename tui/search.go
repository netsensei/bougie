package tui

import (
	"net/url"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/tui/constants"
)

type searchDialogCmpnt int

const (
	searchIn searchDialogCmpnt = iota
	cancelBtn
	okBtn
)

type Search struct {
	searchIn    textinput.Model
	url         string
	Width       int
	Height      int
	activeCmpnt searchDialogCmpnt
}

func NewSearch(Width int, Height int) Search {
	return Search{
		Width:       Width,
		Height:      Height,
		activeCmpnt: searchIn,
	}
}

func (m Search) Init() tea.Cmd {
	return nil
}

func (m Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case SearchMsg:
		in := textinput.New()
		in.Prompt = "Search > "
		in.Placeholder = "go to..."
		in.CharLimit = 250
		in.Focus()

		m.searchIn = in
		m.url = msg.url
		m.activeCmpnt = searchIn

		cmd = SetBrowserModeCmd(input)
		return m, cmd

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Enter) {
			value := m.searchIn.Value()

			if value != "" && m.activeCmpnt == okBtn {
				purl, _ := url.Parse(m.url)
				q := purl.Query()
				q.Set("q", value)
				purl.RawQuery = q.Encode()

				cmds = append(cmds, AddHistoryCmd(purl.String()))
				cmds = append(cmds, StartQueryCmd(purl.String()))
				return m, tea.Batch(cmds...)
			}

			if m.activeCmpnt == cancelBtn {
				cmds = append(cmds, CancelSearchCmd())
				return m, tea.Batch(cmds...)
			}
		}

		if key.Matches(msg, constants.Keymap.Tab) {
			switch m.activeCmpnt {
			case searchIn:
				m.searchIn.Blur()
				m.activeCmpnt = okBtn
				cmds = append(cmds, SetBrowserModeCmd(search))
			case okBtn:
				m.searchIn.Blur()
				m.activeCmpnt = cancelBtn
				cmds = append(cmds, SetBrowserModeCmd(search))
			case cancelBtn:
				m.searchIn.Focus()
				m.activeCmpnt = searchIn
				cmds = append(cmds, SetBrowserModeCmd(input))
			}
		}

		m.searchIn, cmd = m.searchIn.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
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
	switch m.activeCmpnt {
	case searchIn:
		okButton = buttonStyle.MarginRight(2).Render("Search")
		cancelButton = buttonStyle.Render("Cancel")
	case cancelBtn:
		okButton = buttonStyle.MarginRight(2).Render("Search")
		cancelButton = activeButtonStyle.Render("Cancel")
	case okBtn:
		okButton = activeButtonStyle.MarginRight(2).Render("Search")
		cancelButton = buttonStyle.Render("Cancel")
	}

	question := lipgloss.NewStyle().Width(100).Align(lipgloss.Center).Render(m.searchIn.View())
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
