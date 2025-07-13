package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/gopher"
)

type QueryMsg struct {
	url string
}

type AddHistoryMsg struct {
	url string
}

type ReadyMsg struct {
	url     string
	doc     string
	content string
	links   []map[int]string
}

type RedrawMsg struct {
	content  string
	position int
}

type ModeMsg mode

func SetBrowserModeCmd(mode mode) tea.Cmd {
	return func() tea.Msg {
		return ModeMsg(mode)
	}
}

func StartQueryCmd(url string) tea.Cmd {
	return func() tea.Msg {
		return QueryMsg{
			url: url,
		}
	}
}

func AddHistoryCmd(url string) tea.Cmd {
	return func() tea.Msg {
		return AddHistoryMsg{
			url: url,
		}
	}
}

func SendQueryCmd(url string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.TODO()
		request := gopher.New(url)
		response, _ := request.Do(ctx)

		content := ""
		links := []map[int]string{}

		// Parse the response
		switch request.ItemType {
		case gopher.ItemTypeText:
			content = response.Body
		case gopher.ItemTypeDirectory:
			content, links, _ = gopher.ParseDirectory([]byte(response.Body), 0)
		}

		return ReadyMsg{
			url:     url,
			content: content,
			doc:     response.Body,
			links:   links,
		}
	}
}

func RedrawCmd(doc string, position int) tea.Cmd {
	return func() tea.Msg {
		content, _, _ := gopher.ParseDirectory([]byte(doc), position)
		return RedrawMsg{
			content:  content,
			position: position,
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
