package tui

import (
	"context"
	"fmt"
	purl "net/url"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/gopher"
)

type GopherDocumentQueryMsg struct {
	url     string
	request *gopher.Request
}

type GopherFileQueryCmd struct {
	url     string
	request *gopher.Request
}

type AddHistoryMsg struct {
	url string
}

type SearchMsg struct {
	url string
}

type CancelSearchMsg struct{}

type ReadyMsg struct {
	url     string
	doc     string
	content string
	links   []map[int]string
}

type FileSavedMsg struct {
	url      string
	resource string
}

type ErrorMsg struct {
	url string
	err error
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

			switch request.ItemType {
			case gopher.ItemTypeBinary:
				fallthrough
			case gopher.ItemTypeDOS:
				fallthrough
			case gopher.ItemTypeGIF:
				fallthrough
			case gopher.ItemTypeHex:
				fallthrough
			case gopher.ItemTypeImage:
				return GopherFileQueryCmd{
					url:     url,
					request: request,
				}
			}

			if request.ItemType == gopher.ItemTypeSEA {
				if len(purl.Query()) == 0 {
					return SearchMsg{
						url: url,
					}
				}
			}

			return GopherDocumentQueryMsg{
				url:     url,
				request: request,
			}
		}

		return ErrorMsg{
			url: url,
			err: fmt.Errorf("unrecognized scheme: %v", purl.Scheme),
		}
	}
}

func FetchDocumentGopherCmd(request *gopher.Request, url string) tea.Cmd {
	return func() tea.Msg {
		var content string
		var links []map[int]string
		var err error

		ctx := context.TODO()
		response, err := request.Do(ctx)
		if err != nil {
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("could not reach gopher server: %v", err),
			}
		}

		// Process the response
		switch request.ItemType {
		case gopher.ItemTypeText:
			content = string(response.Body)
		case gopher.ItemTypeSEA:
			fallthrough
		case gopher.ItemTypeDirectory:
			content, links, _ = gopher.ParseDirectory(response.Body, 0)
		}

		return ReadyMsg{
			url:     url,
			content: content,
			doc:     string(response.Body),
			links:   links,
		}
	}
}

func SaveFileGopherCmd(request *gopher.Request, url string) tea.Cmd {
	return func() tea.Msg {
		resource := filepath.Base(request.Selector)

		ctx := context.TODO()
		response, err := request.Do(ctx)
		if err != nil {
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("could not reach gopher server: %v", err),
			}
		}

		filePath := filepath.Join(config.DownloadsDir, resource)

		err = os.WriteFile(filePath, response.Body, 0644)
		if err != nil {
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("could not save file: %v", err),
			}
		}

		return FileSavedMsg{
			url:      url,
			resource: resource,
		}
	}
}

func CancelSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return CancelSearchMsg{}
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
