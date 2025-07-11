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
	content string
	links   map[int]string
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

		content := ""
		links := make(map[int]string)

		// Parse the response
		switch request.ItemType {
		case gopher.ItemTypeText:
			content = response.Body
		case gopher.ItemTypeDirectory:
			content, links, _ = gopher.ParseDirectory([]byte(response.Body))
		}

		return ReadyMsg{
			content: content,
			links:   links,
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
