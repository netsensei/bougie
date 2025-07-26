package tui

import (
	"context"
	"fmt"
	purl "net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/gopher"
)

type GopherQueryMsg struct {
	url     string
	request *gopher.Request
}

type SchemeErrorMsg struct{}

type AddHistoryMsg struct {
	url string
}

type SearchMsg struct {
	url string
}

type ReadyMsg struct {
	url     string
	doc     string
	content string
	links   []map[int]string
	err     error
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
		purl, _ := purl.Parse(url)

		switch purl.Scheme {
		case "gopher":
			request := gopher.New(purl.String())

			if request.ItemType == gopher.ItemTypeSEA {
				if len(purl.Query()) == 0 {
					return SearchMsg{
						url: url,
					}
				}
			}

			return GopherQueryMsg{
				url:     url,
				request: request,
			}
		}

		return SchemeErrorMsg{}
	}
}

func SendGopherQueryCmd(request *gopher.Request, url string) tea.Cmd {
	return func() tea.Msg {
		content := ""
		doc := ""
		var links []map[int]string
		var err error

		content, doc, links, err = handleGopher(request)
		if err != nil {
			return ReadyMsg{
				url:     url,
				content: "",
				doc:     "",
				links:   links,
				err:     err,
			}
		}

		return ReadyMsg{
			url:     url,
			content: content,
			doc:     doc,
			links:   links,
			err:     nil,
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

func RedrawCmd(doc string, position int) tea.Cmd {
	return func() tea.Msg {
		content, _, _ := gopher.ParseDirectory([]byte(doc), position)
		return RedrawMsg{
			content:  content,
			position: position,
		}
	}
}

func handleGopher(request *gopher.Request) (string, string, []map[int]string, error) {
	var links []map[int]string
	var content string

	ctx := context.TODO()
	response, err := request.Do(ctx)
	if err != nil {
		err = fmt.Errorf("could not reach gopher server: %v", err)
		return "", "", nil, err
	}

	// Process the response
	switch request.ItemType {
	case gopher.ItemTypeText:
		content = response.Body
	case gopher.ItemTypeSEA:
		fallthrough
	case gopher.ItemTypeDirectory:
		content, links, _ = gopher.ParseDirectory([]byte(response.Body), 0)
	}

	return content, response.Body, links, nil
}
