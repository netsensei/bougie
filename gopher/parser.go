package gopher

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Item interface {
	Parse(raw []byte)
	Render() string
}

type Text struct {
	Content string
}

func (t *Text) Parse(raw []byte) {
	t.Content = string(raw)
}

func (t *Text) Render() string {
	return t.Content
}

type Menu struct {
	Lines  []*Line
	Active int
}

type Link struct {
	URL      url.URL
	Position int
}

func NewMenu() *Menu {
	return &Menu{
		Active: 0,
	}
}

func (m *Menu) Links() []Link {
	res := make([]Link, 0)
	for k, l := range m.Lines {
		switch l.Type {
		case "0", "1", "3", "7":
			url := url.URL{
				Scheme: "gopher",
				Host:   l.Host,
				Path:   l.Path,
			}

			res = append(res, Link{
				URL:      url,
				Position: k,
			})
		}
	}

	return res
}

func (m *Menu) ActiveLink() *Link {
	l := m.Links()
	return &l[m.Active]
}

type Line struct {
	Type string
	Text string
	Path string
	Host string
	Port string
}

// Should be part of Menu??
func Parse(raw []byte) Menu {
	p := Menu{
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
