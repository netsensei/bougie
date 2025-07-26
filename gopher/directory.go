package gopher

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/bougie/tui/constants"
)

func ParseDirectory(body []byte, active int) (string, []map[int]string, error) {
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

	documentStyle := lipgloss.NewStyle()
	//	Background(lipgloss.Color("#7D56F4"))

	textStyle := lipgloss.NewStyle().
		Inherit(documentStyle).
		Width(constants.WindowWidth).
		Foreground(lipgloss.Color("#FAFAFA"))

	typeStyle := lipgloss.NewStyle().
		Inherit(documentStyle).
		Width(6)

	linkStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	activeLinkStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Bold(true).
		Foreground(lipgloss.Color("#CC56F4"))

	// Let's go
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)

	var links []map[int]string

	doc := strings.Builder{}

	lnumber := 0
	for scanner.Scan() {
		var line string

		st := scanner.Text()

		if len(st) == 0 {
			continue
		}

		lp := strings.Split(st[1:], "\t")
		if len(lp) < 4 {
			continue
		}

		itype := st[:1]

		switch itype {
		case ItemTypeNCInformation:
			text := textStyle.Render(lp[0])
			itemType := typeStyle.Render("")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		case ItemTypeText:
			fallthrough
		case ItemTypeSEA:
			fallthrough
		case ItemTypeDirectory:
			text := textStyle.Render(lp[0])
			itemType := linkStyle.Render(types[itype])
			if lnumber == active || active == 0 {
				itemType = activeLinkStyle.Render(types[itype])
				active = -1 // reset active to -1 so we don't highlight again
			}

			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

			host := lp[2]
			if lp[3] != "" {
				host += ":" + lp[3]
			}

			url := url.URL{
				Scheme: "gopher",
				Host:   host,
				Path:   itype + lp[1],
			}

			links = append(links, map[int]string{lnumber: url.String()})

		case ItemTypeCCSO:
			fallthrough
		case ItemType3270:
			fallthrough
		case ItemTypeHex:
			fallthrough
		case ItemTypeDOS:
			fallthrough
		case ItemTypeUUE:
			fallthrough
		case ItemTypeTelnet:
			fallthrough
		case ItemTypeBinary:
			fallthrough
		case ItemTypeAlt:
			fallthrough
		case ItemTypeGIF:
			fallthrough
		case ItemTypeImage:
			text := textStyle.Render(lp[0])
			itemType := linkStyle.Render(types[itype])
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

			// line = linkStyle.Render(types[itype]) + " " + textStyle.Render(lp[0]) + "\t" + textStyle.Render(lp[1])

		default:
			text := textStyle.Render(lp[0])
			itemType := typeStyle.Render("[***]")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)
		}

		lnumber++

		doc.WriteString(line + "\n")
	}

	return doc.String(), links, nil
}
