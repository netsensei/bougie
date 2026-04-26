package gopher

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"charm.land/lipgloss/v2"
)

var types = map[string]string{
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

type Item struct {
	ItemType   string
	Display    string
	Selector   string
	Host       string
	Port       string
	LineNumber int
}

func (i Item) Url() string {
	switch i.ItemType {
	case ItemTypeNCInformation, ItemTypeError:
		return ""

	// HTML / external URL (http, gemini, ftp, etc)
	case ItemtypeNCHTML:
		if strings.HasPrefix(i.Selector, "URL:") {
			return strings.TrimPrefix(i.Selector, "URL:")
		}

	case ItemTypeTelnet, ItemType3270:
		host := i.Host
		if i.Port != "" {
			host += ":" + i.Port
		}
		return "telnet://" + host + "/"
	}

	selector := i.Selector
	if selector == "" {
		selector = "/"
	} else if !strings.HasPrefix(selector, "/") {
		selector = "/" + selector
	}

	host := i.Host
	if i.Port != "" && i.Port != "70" {
		host += ":" + i.Port
	}

	u := &url.URL{
		Scheme: "gopher",
		Host:   host,
		Path:   "/" + i.ItemType + selector,
	}

	return u.String()
}

type Directory struct {
	Items []Item
}

func Parse(body []byte) (Directory, error) {
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)

	var items []Item

	lnumber := 0
	for scanner.Scan() {
		var item Item

		st := strings.TrimRight(scanner.Text(), "\r")

		if len(st) == 0 {
			continue
		}

		if st == "." {
			break
		}

		lp := strings.SplitN(st, "\t", 5)
		if len(lp) < 4 {
			continue
		}

		itemType := lp[0][:1]
		display := ""
		if len(lp[0]) > 1 {
			display = lp[0][1:]
		}

		item = Item{
			ItemType:   itemType,
			Display:    display,
			Selector:   lp[1],
			Host:       lp[2],
			Port:       lp[3],
			LineNumber: lnumber,
		}

		items = append(items, item)
		lnumber++
	}

	if err := scanner.Err(); err != nil {
		return Directory{}, err
	}

	return Directory{Items: items}, nil
}

func (d Directory) FirstLink() int {
	for _, item := range d.Items {
		if item.Url() != "" {
			return item.LineNumber
		}
	}
	return -1
}

func (d Directory) Links() []map[int]string {
	var links []map[int]string

	for _, item := range d.Items {
		url := item.Url()
		if url != "" {
			links = append(links, map[int]string{item.LineNumber: url})
		}
	}
	return links
}

func (d Directory) Render(active int) string {
	var sb strings.Builder

	for _, item := range d.Items {
		var line string

		switch item.ItemType {
		case ItemTypeNCInformation:
			text := textStyle.Render(item.Display)
			itemType := typeStyle.Render("")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		case ItemTypeText, ItemTypeDirectory, ItemTypeHex, ItemTypeDOS,
			ItemTypeUUE, ItemTypeTelnet, ItemTypeBinary, ItemTypeAlt,
			ItemTypeGIF, ItemTypeImage, ItemTypeSEA:
			text := textStyle.Render(item.Display)
			itemType := linkStyle.Render(types[item.ItemType])
			if item.LineNumber == active {
				itemType = activeLinkStyle.Render(types[item.ItemType])
			}

			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		case ItemType3270, ItemTypeCCSO:
			text := textStyle.Render(item.Display)
			itemType := linkStyle.Render(types[item.ItemType])
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)

		default:
			text := textStyle.Render(item.Display)
			itemType := typeStyle.Render("[***]")
			line = lipgloss.JoinHorizontal(lipgloss.Top, itemType, text)
		}

		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}
