package gopher

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	styleLink := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	styleActiveLink := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CC56F4"))

	styleText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	// Let's go
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer

	var links []map[int]string

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
			line = "\t" + styleText.Render(lp[0]) + "\n"

		case ItemTypeText:
		case ItemTypeDirectory:
			line = styleLink.Render(types[itype]) + "\t" + styleText.Render(lp[0]) + "\n"
			if lnumber == active || active == 0 {
				line = styleActiveLink.Render(types[itype]) + "\t" + styleText.Render(lp[0]) + "\n"
				active = -1 // reset active to -1 so we don't highlight again
			}

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
		case ItemTypeSEA:
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
			line = styleLink.Render(types[itype]) + "\t" + styleText.Render(lp[0]) + "\t" + styleText.Render(lp[1]) + "\n"

		default:
			line = styleText.Render("[***]") + "\t" + styleText.Render(lp[0]) + "\n"

		}

		lnumber++

		buffer.WriteString(line)
	}

	return buffer.String(), links, nil
}
