package tui

import (
	"context"
	"fmt"
	purl "net/url"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/gemini"
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

type GeminiQueryMsg struct {
	url     string
	request *gemini.Request
}

type FetchGemTextGeminiMsg struct {
	capsule *gemini.Capsule
	url     string
}

type SaveFileGeminiMsg struct {
	capsule *gemini.Capsule
	url     string
}

type RedirectQueryGeminiMsg struct {
	url string
}

type AddHistoryMsg struct {
	url string
}

type SearchMsg struct {
	url string
}

type CancelSearchMsg struct{}

type ReadyMsg struct {
	currentUrl string
	doc        string
	content    string
	scheme     string
	links      []map[int]string
}

type ViewSourceMsg struct {
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
	content string
	active  int
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

		case "gemini":
			request, _ := gemini.NewRequest(purl.String())
			return GeminiQueryMsg{
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

func FetchCapsuleGeminiCmd(request *gemini.Request, url string) tea.Cmd {
	return func() tea.Msg {
		capsule, _ := gemini.FetchCapsuleGemini(request)

		switch capsule.Status {
		case 1:
			// Input required
		case 2:
			// Successful response
			if len(capsule.Body) == 0 {
				return ErrorMsg{
					url: url,
					err: fmt.Errorf("received empty document"),
				}
			}

			mimetype := strings.Split(capsule.Information, ";")[0]
			if mimetype == "" {
				mimetype = "text/gemini"
			}

			typeSubtype := strings.Split(mimetype, "/")
			if len(typeSubtype) != 2 {
				return ErrorMsg{
					url: url,
					err: fmt.Errorf("invalid mime type received: %v", mimetype),
				}
			}

			if typeSubtype[0] == "text" && typeSubtype[1] == "gemini" {
				return FetchGemTextGeminiMsg{
					capsule: &capsule,
					url:     url,
				}
			} else {
				return SaveFileGeminiMsg{
					capsule: &capsule,
					url:     url,
				}
			}

		case 3:
			if len(capsule.Information) == 0 {
				return ErrorMsg{
					url: url,
					err: fmt.Errorf("redirection with no target"),
				}
			}

			if strings.Index(capsule.Information, "gemini://") == 0 {
				base, err := purl.Parse(url)
				if err != nil {
					return ErrorMsg{
						url: url,
						err: fmt.Errorf("invalid base URL for redirection: %v", err),
					}
				}
				rel, err := purl.Parse(capsule.Information)
				if err != nil {
					return ErrorMsg{
						url: url,
						err: fmt.Errorf("invalid redirection URL: %v", err),
					}
				}
				newUrl := base.ResolveReference(rel).String()

				return RedirectQueryGeminiMsg{
					url: newUrl,
				}
			}

		case 4:
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("temporary failure: %v", capsule.Information),
			}
		case 5:
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("permanent failure: %v", capsule.Information),
			}
		case 6:
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("client certificate required (not supported): %v", capsule.Information),
			}
		default:
			return ErrorMsg{
				url: url,
				err: fmt.Errorf("invalid status code received: %v", capsule.Status),
			}
		}

		return nil
	}
}

func FetchGemTextGeminiCmd(capsule *gemini.Capsule, currentUrl string) tea.Cmd {
	return func() tea.Msg {
		// content := ""
		// links := []map[int]string{}

		content, links, _ := gemini.ParseGemText(capsule.Body, currentUrl, 0)
		// content = "Bougie, a tiny sparking Gopher browser"
		links = append(links, map[int]string{1: "gemini://example.com"})
		return ReadyMsg{
			currentUrl: currentUrl,
			content:    content,
			doc:        string(capsule.Body),
			links:      links,
			scheme:     "gemini",
		}
	}
}

func SaveFileGeminiCmd(capsule *gemini.Capsule, url string) tea.Cmd {
	return func() tea.Msg {
		purl, _ := purl.Parse(url)
		resource := filepath.Base(purl.Path)

		if resource == "" || resource == "/" || resource == "." {
			resource = "index.gmi"
		}

		filePath := filepath.Join(config.DownloadsDir, resource)

		err := os.WriteFile(filePath, capsule.Body, 0644)
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

func FetchDocumentGopherCmd(request *gopher.Request, currentUrl string) tea.Cmd {
	return func() tea.Msg {
		var content string
		var links []map[int]string
		var err error

		ctx := context.TODO()
		response, err := request.Do(ctx)
		if err != nil {
			return ErrorMsg{
				url: currentUrl,
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
			currentUrl: currentUrl,
			content:    content,
			doc:        string(response.Body),
			links:      links,
			scheme:     "gopher",
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

func ViewSourceCmd() tea.Cmd {
	return func() tea.Msg {
		return ViewSourceMsg{}
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

func RedrawCmd(scheme string, currentUrl string, doc string, active int) tea.Cmd {
	return func() tea.Msg {
		var content string

		if scheme == "gemini" {
			content, _, _ = gemini.ParseGemText([]byte(doc), currentUrl, active)
		} else {
			content, _, _ = gopher.ParseDirectory([]byte(doc), active)
		}

		return RedrawMsg{
			content: content,
			active:  active,
		}
	}
}
