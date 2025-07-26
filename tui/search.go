package tui

import (
	"net/url"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/tui/constants"
)

type Search struct {
	input textinput.Model
	mode  mode
	url   string
}

func NewSearch() Search {
	return Search{}
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

				if value != "" {
					purl, _ := url.Parse(m.url)
					q := purl.Query()
					q.Set("q", value)
					purl.RawQuery = q.Encode()

					cmds = append(cmds, AddHistoryCmd(purl.String()))
					cmds = append(cmds, StartQueryCmd(purl.String()))
					return m, tea.Batch(cmds...)
				}
			}

			m.input, cmd = m.input.Update(msg)
		}

		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Search) View() string {
	return m.input.View()
}
