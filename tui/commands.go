package tui

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type QueryMsg struct {
	url string
}

type ReadyMsg string

type ModeMsg mode

func setBrowserMode(mode mode) tea.Cmd {
	return func() tea.Msg {
		return mode
	}
}

func queryCmd(url string) tea.Cmd {
	// Parse the Url
	return func() tea.Msg {
		return QueryMsg{
			url: url,
		}
	}
}

// func doRequestCmd() tea.Cmd {
// 	// Return a request
// }

// func parseResponse() tea.Cmd {
// 	// Return a parsed response
// }

func GetContent(foo string) tea.Cmd {
	return func() tea.Msg {
		content, _ := os.ReadFile("bacon.txt")
		time.Sleep(3 * time.Second)
		return ReadyMsg(content)
	}
}
