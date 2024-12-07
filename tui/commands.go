package tui

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type LoadingMsg struct {
	url string
}

type ReadyMsg string

func load(url string) tea.Cmd {
	return func() tea.Msg {
		return LoadingMsg{
			url: url,
		}
	}
}

func GetContent() tea.Msg {
	content, _ := os.ReadFile("bacon.txt")
	time.Sleep(3 * time.Second)
	return ReadyMsg(content)
}
