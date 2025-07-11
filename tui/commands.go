package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/gopher"
)

type QueryMsg struct {
	url string
}

type ReadyMsg struct {
	response *gopher.Response
}

type ModeMsg mode

func setBrowserModeCmd(mode mode) tea.Cmd {
	return func() tea.Msg {
		return ModeMsg(mode)
	}
}

func SendQueryCmd(url string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		request := gopher.New(url)
		response, _ := request.Do(ctx)

		// Parse the response

		return ReadyMsg{
			response: response,
		}
	}
}

// func GetContent(foo string) tea.Cmd {
// 	return func() tea.Msg {
// 		content, _ := os.ReadFile("bacon.txt")
// 		time.Sleep(2 * time.Second)
// 		return ReadyMsg(content)
// 	}
// }
