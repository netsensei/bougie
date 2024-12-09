package gopher

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Text struct {
	Content string
}

func NewText(raw []byte) *Text {
	return &Text{
		Content: string(raw),
	}
}

func (t *Text) Render() string {
	return t.Content
}

type Menu struct {
	Lines   []*Line
	Links   []*Link
	Active  int
	Content string
}

func NewMenu(raw []byte) *Menu {
	m := &Menu{
		Active: 0,
	}

	reader := bytes.NewReader(raw)
	scanner := bufio.NewScanner(reader)

	pos := 0

	for scanner.Scan() {
		l := &Line{}
		st := scanner.Text()

		if len(st) < 1 {
			continue
		}

		lp := strings.Split(st[1:], "\t")
		if len(lp) < 4 {
			continue
		}

		l.Type = st[:1]

		switch l.Type {
		case "i":
			l.Text = lp[0]
		case "0", "1", "3", "7":
			path := l.Type
			if path != "" {
				path = path + lp[1]
			}

			host := lp[2]
			if lp[3] != "" {
				host = lp[2] + ":" + lp[3]
			}

			l.Path = path
			l.Host = host
			l.Text = lp[0]

			url := url.URL{
				Scheme: "gopher",
				Host:   host,
				Path:   path,
			}

			m.Links = append(m.Links, &Link{
				URL:      url.String(),
				Position: pos,
			})

			pos++
		}

		m.Lines = append(m.Lines, l)

	}

	return m
}

func ParseMenu(raw []byte) *Menu {
	p := &Menu{
		Active: 0,
	}

	reader := bytes.NewReader(raw)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		l := &Line{}

		st := scanner.Text()

		if len(st) < 1 {
			continue
		}

		lp := strings.Split(st[1:], "\t")
		if len(lp) < 4 {
			continue
		}

		l.Type = st[:1]

		switch l.Type {
		case "i":
			l.Text = lp[0]
		case "0", "1", "3", "7":
			path := l.Type
			if path != "" {
				path = path + lp[1]
			}

			host := lp[2]
			if lp[3] != "" {
				host = lp[2] + ":" + lp[3]
			}

			l.Path = path
			l.Host = host
			l.Text = lp[0]
		}

		p.Lines = append(p.Lines, l)
	}

	return p
}

func (m *Menu) Render() string {
	types := map[string]string{
		"0": "[TXT]",
		"1": "[SUB]",
		"2": "[CCS]",
		"3": "[ERR]",
		"4": "[HEX]",
		"5": "[DOS]",
		"6": "[UUE]",
		"7": "[SEA]",
		"8": "[TEL]",
		"9": "[BIN]",
		"+": "[ALT]",
		"g": "[GIF]",
		"I": "[IMG]",
		"T": "[327]",
	}

	styleType := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	styleActiveLink := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CC56F4"))

	styleText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	res := ""

	pLink := 0

	for _, l := range m.Lines {
		var tmp string
		switch l.Type {
		case "i":
			tmp = "\t" + styleText.Render(l.Text) + "\n"
		case "0", "1", "3", "7":
			if m.Active == pLink {
				tmp = styleActiveLink.Render(types[l.Type]) + "\t" + styleText.Render(l.Text) + "\n"
			} else {
				tmp = styleType.Render(types[l.Type]) + "\t" + styleText.Render(l.Text) + "\n"
			}

			pLink++
		default:
			tmp = styleType.Render("[***]") + "\t" + styleText.Render(l.Text) + "\n"
		}

		res = res + tmp
	}

	return res
}

type Line struct {
	Type string
	Text string
	Path string
	Host string
	Port string
}

type Link struct {
	URL      string
	Position int
}
