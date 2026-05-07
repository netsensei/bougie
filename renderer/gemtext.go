package renderer

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/netsensei/bougie/gemini"
	"github.com/netsensei/bougie/theme"
)

var documentStyle = lipgloss.NewStyle()

var typeStyle = lipgloss.NewStyle().
	Inherit(documentStyle)

var textStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Foreground(theme.Default.ContentMuted).
	// MarginTop(1).
	MarginBottom(1)

var headingStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Bold(true).
	Foreground(theme.Default.Heading).
	MarginBottom(1)

var linkStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Bold(true).
	Foreground(theme.Default.Link)

var activeLinkStyle = lipgloss.NewStyle().
	Inherit(typeStyle).
	Bold(true).
	Foreground(theme.Default.LinkActive)

func RenderGemText(g gemini.GemText, active int) string {
	var sb strings.Builder

	spacer := typeStyle.Render(strings.Repeat(" ", 3))

	var wrapped, out string
	for i, node := range g.Nodes {
		switch node.Type {
		case gemini.TextNodeType:
			wrapped = wrapText(node.Text, 80)
			out = textStyle.Render(wrapped)
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, out) + "\n")

		case gemini.HeadingNodeType:
			wrapped = wrapText(strings.Repeat("#", node.Level)+" "+node.Text, 80)
			out = headingStyle.Render(wrapped)
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, out) + "\n")

		case gemini.LinkNodeType:
			wrapped = wrapText(node.Alt, 80)
			if active == node.LineNumber {
				out = activeLinkStyle.Render(wrapped)
			} else {
				out = linkStyle.Render(wrapped)
			}

			// Peek ahead
			nl := "\n"
			if i < len(g.Nodes)-1 {
				if g.Nodes[i+1].Type != gemini.LinkNodeType {
					nl = "\n\n"
				}
			}

			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, out) + nl)

		case gemini.PreformattedNodeType:
			for _, line := range node.Lines {
				sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, line) + "\n")
			}

		case gemini.ListNodeType:
			out := "* " + node.Text

			// Peek ahead
			nl := "\n"
			if i < len(g.Nodes)-1 {
				if g.Nodes[i+1].Type != gemini.ListNodeType {
					nl = "\n\n"
				}
			}
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, out) + nl)

		case gemini.BlockquoteNodeType:
			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, spacer, "> "+node.Text) + "\n")
		}
	}

	return sb.String()
}

func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var sb strings.Builder
	words := strings.Fields(text)
	lineLength := 0

	for i, word := range words {
		wordLength := len(word)
		if i == 0 {
			sb.WriteString(word)
			lineLength = wordLength
		} else if (lineLength + 1 + wordLength) > width {
			sb.WriteString("\n" + word)
			lineLength = wordLength
		} else {
			sb.WriteString(" " + word)
			lineLength += 1 + wordLength
		}
	}

	return sb.String()
}
