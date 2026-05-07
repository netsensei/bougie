package gemini

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"
)

type NodeType int

const (
	TextNodeType NodeType = iota
	LinkNodeType
	PreformattedNodeType
	HeadingNodeType
	ListNodeType
	BlockquoteNodeType
)

type Scope int

const (
	ScopeInternal Scope = iota
	ScopeExternal
)

type GemText struct {
	Nodes []Node
}

func (g GemText) Links(currentUrl string) []map[int]string {
	var links []map[int]string

	for _, node := range g.Nodes {
		if node.Type == LinkNodeType && node.URL != "" {
			node.Scope = ScopeInternal
			if !strings.Contains(node.URL, "://") {
				base, err := url.Parse(currentUrl)
				if err != nil {
					continue
				}

				u, err := url.Parse(node.URL)
				if err != nil {
					continue
				}

				node.URL = base.ResolveReference(u).String()
				node.Scope = ScopeExternal
			}

			links = append(links, map[int]string{node.LineNumber: node.URL})
		}
	}

	return links
}

func (g GemText) FirstLink() int {
	for _, node := range g.Nodes {
		if node.Type == LinkNodeType && node.URL != "" {
			return node.LineNumber
		}
	}
	return -1
}

type Node struct {
	Type       NodeType
	Text       string
	Level      int
	Alt        string
	URL        string
	Scope      Scope
	Lines      []string
	LineNumber int
}

func Parse(body []byte) (GemText, error) {
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)

	var nodes []Node
	var paragraph []string

	inPre := false
	preLines := []string{}

	flushParagraph := func() {
		if len(paragraph) > 0 {
			nodes = append(nodes, Node{
				Type: TextNodeType,
				Text: strings.Join(paragraph, "\n"),
			})
			paragraph = nil
		}
	}

	lnumber := 0
	for scanner.Scan() {
		lnumber++
		st := strings.TrimRight(scanner.Text(), "\r")

		// -- preformatted block
		if strings.HasPrefix(st, "```") {
			flushParagraph()
			if inPre {
				nodes = append(nodes, Node{
					Type:  PreformattedNodeType,
					Lines: append([]string{}, st),
				})
				preLines = nil
				inPre = false
			} else {
				inPre = true
			}
			continue
		}

		if inPre {
			preLines = append(preLines, st)
			continue
		}

		// -- blank line
		if len(st) == 0 {
			// flush paragraph
			continue
		}

		// -- Headings
		if strings.HasPrefix(st, "#") {
			// flush paragraph
			flushParagraph()
			nodes = append(nodes, Node{
				Type:       HeadingNodeType,
				Text:       strings.TrimLeft(st, "# "),
				Level:      1,
				LineNumber: lnumber,
			})
			continue
		}

		if strings.HasPrefix(st, "##") {
			// flush paragraph
			flushParagraph()
			nodes = append(nodes, Node{
				Type:       HeadingNodeType,
				Text:       strings.TrimLeft(st, "##"),
				Level:      2,
				LineNumber: lnumber,
			})
			continue
		}

		if strings.HasPrefix(st, "###") {
			// flush paragraph
			flushParagraph()
			nodes = append(nodes, Node{
				Type:       HeadingNodeType,
				Text:       strings.TrimLeft(st, "###"),
				Level:      3,
				LineNumber: lnumber,
			})
			continue
		}

		// -- Links
		if strings.HasPrefix(st, "=>") {
			// flush paragraph
			flushParagraph()
			link := strings.TrimSpace(st[2:])
			parts := strings.Fields(link)

			var url, alt string
			if len(parts) > 0 {
				url = parts[0]
			}

			if len(parts) > 1 {
				alt = strings.Join(parts[1:], " ")
			}

			nodes = append(nodes, Node{
				Type:       LinkNodeType,
				URL:        url,
				Alt:        alt,
				LineNumber: lnumber,
			})
			continue
		}

		// -- list items
		if strings.HasPrefix(st, "*") {
			// flush paragraph
			flushParagraph()
			nodes = append(nodes, Node{
				Type:       ListNodeType,
				Text:       strings.TrimLeft(st, "* "),
				LineNumber: lnumber,
			})
			continue
		}

		paragraph = append(paragraph, st)
	}

	// flush remaining paragraph
	flushParagraph()

	return GemText{Nodes: nodes}, nil
}
