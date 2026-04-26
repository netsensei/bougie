package renderer

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/netsensei/bougie/gopher"
	"github.com/netsensei/bougie/theme"
	"github.com/netsensei/bougie/tui/constants"
)

var gopherTypes = map[string]string{
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

var DocumentStyle = lipgloss.NewStyle()

var TextStyle = lipgloss.NewStyle().
	Inherit(DocumentStyle).
	Width(constants.WindowWidth).
	Foreground(theme.Default.ContentFg)

var TypeStyle = lipgloss.NewStyle().
	Inherit(DocumentStyle).
	Width(6)

var LinkStyle = lipgloss.NewStyle().
	Inherit(TypeStyle).
	Bold(true).
	Foreground(theme.Default.Link)

var gopherActiveLinkStyle = lipgloss.NewStyle().
	Inherit(TypeStyle).
	Bold(true).
	Foreground(theme.Default.LinkActive)

func RenderGopherDirectory(d gopher.Directory, active int) string {
	var sb strings.Builder

	for _, item := range d.Items {
		var line string

		switch item.ItemType {
		case gopher.ItemTypeNCInformation:
			text := TextStyle.Render(item.Display)
			itemType := TypeStyle.Render("")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		case gopher.ItemTypeText, gopher.ItemTypeDirectory, gopher.ItemTypeHex, gopher.ItemTypeDOS,
			gopher.ItemTypeUUE, gopher.ItemTypeTelnet, gopher.ItemTypeBinary, gopher.ItemTypeAlt,
			gopher.ItemTypeGIF, gopher.ItemTypeImage, gopher.ItemTypeSEA:
			text := TextStyle.Render(item.Display)
			itemType := LinkStyle.Render(gopherTypes[item.ItemType])
			if item.LineNumber == active {
				itemType = gopherActiveLinkStyle.Render(gopherTypes[item.ItemType])
			}
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		case gopher.ItemType3270, gopher.ItemTypeCCSO:
			text := TextStyle.Render(item.Display)
			itemType := LinkStyle.Render(gopherTypes[item.ItemType])
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		default:
			text := TextStyle.Render(item.Display)
			itemType := TypeStyle.Render("[***]")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)
		}

		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}
