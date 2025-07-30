package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/history"
	"github.com/netsensei/bougie/tui/constants"
	"github.com/spf13/viper"
)

type mode int

const (
	nav mode = iota
	view
	search
	input
	save
)

type status int

const (
	ready status = iota
	loading
	saving
	saved
	errored
)

type Browser struct {
	status     Status
	canvas     Canvas
	search     Search
	navigation Navigation
	history    *history.History
	mode       mode
	quitting   bool
	ready      bool
}

func initBrowser() (tea.Model, tea.Cmd) {
	status := NewStatus()
	navigation := NewNavigation()
	canvas := NewCanvas()

	m := Browser{
		// mode:       nav,
		status:     status,
		navigation: navigation,
		canvas:     canvas,
		history: &history.History{
			Position: 0,
			Length:   0,
		},
	}

	return m, func() tea.Msg { return nil }
}

func (m Browser) Init() tea.Cmd {
	var cmds []tea.Cmd

	home := viper.GetString("general.home")

	cmds = append(cmds, m.status.Init())

	cmds = append(cmds, SetBrowserModeCmd(view))
	cmds = append(cmds, AddHistoryCmd(home))
	cmds = append(cmds, StartQueryCmd(home))

	return tea.Batch(cmds...)
}

func (m Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.search = NewSearch(constants.WindowWidth, constants.WindowHeight)
			m.ready = true
		}

		if m.ready {
			constants.WindowHeight = msg.Height - lipgloss.Height(m.navigation.View()) - lipgloss.Height(m.status.View())
			constants.WindowWidth = msg.Width

			m.search.Width = constants.WindowWidth
			m.search.Height = constants.WindowHeight
		}

	case ModeMsg:
		m.mode = mode(msg)

	case AddHistoryMsg:
		m.history.Add(msg.url)

	case GopherDocumentQueryMsg:
		cmds = append(cmds, FetchDocumentGopherCmd(msg.request, msg.url))

	case GopherFileQueryCmd:
		cmds = append(cmds, SetBrowserModeCmd(save))
		cmds = append(cmds, SaveFileGopherCmd(msg.request, msg.url))

	case FileSavedMsg:
		cmds = append(cmds, SetBrowserModeCmd(view))

	case ErrorMsg:
		cmds = append(cmds, SetBrowserModeCmd(view))

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.mode {
		case nav:
			m.navigation, cmd = m.navigation.Update(msg)
		case view:
			m.canvas, cmd = m.canvas.Update(msg)
		case search:
			fallthrough
		case input:
			m.search, cmd = m.search.Update(msg)
		}

		if m.mode == view || m.mode == search {
			if key.Matches(msg, constants.Keymap.PageBackward) {
				if m.history.Length > 0 {
					m.history.Backward()
					url := m.history.Current()
					if url != "" {
						cmd = StartQueryCmd(url)
					}
				}
			}

			if key.Matches(msg, constants.Keymap.PageForward) {
				if m.history.Length > 0 {
					m.history.Forward()
					url := m.history.Current()
					if url != "" {
						cmd = StartQueryCmd(url)
					}
				}
			}

			if key.Matches(msg, constants.Keymap.Home) {
				home := viper.GetString("general.home")

				cmds = append(cmds, AddHistoryCmd(home))
				cmds = append(cmds, StartQueryCmd(home))

				return m, tea.Batch(cmds...)
			}

			if key.Matches(msg, constants.Keymap.Reload) {
				url := m.history.Current()
				cmds = append(cmds, StartQueryCmd(url))

				return m, tea.Batch(cmds...)
			}
		}

		return m, cmd
	}

	m.navigation, cmd = m.navigation.Update(msg)
	cmds = append(cmds, cmd)

	m.status, cmd = m.status.Update(msg)
	cmds = append(cmds, cmd)

	m.canvas, cmd = m.canvas.Update(msg)
	cmds = append(cmds, cmd)

	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var (
	InputStyle = lipgloss.NewStyle()
)

func (m Browser) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	navStyle := lipgloss.NewStyle()

	canvasStyle := lipgloss.NewStyle().
		Padding(0, 1)

	searchStyle := lipgloss.NewStyle()

	statusStyle := lipgloss.NewStyle()

	navKey := navStyle.Render(m.navigation.View())
	statusKey := statusStyle.Render(m.status.View())

	if m.mode == view || m.mode == nav || m.mode == save {
		canvasKey := canvasStyle.Render(m.canvas.View())
		return lipgloss.JoinVertical(lipgloss.Top, navKey, canvasKey, statusKey)
	} else {
		searchKey := searchStyle.Render(m.search.View())
		return lipgloss.JoinVertical(lipgloss.Top, navKey, searchKey, statusKey)
	}
}
