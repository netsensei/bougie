package gemini

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Request struct {
	URL         *url.URL
	Certificate *tls.Certificate
}

type Capsule struct {
	Status      int
	Information string
	Body        []byte
}

func NewRequest(raw string) (*Request, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	u.User = nil
	u.Fragment = ""
	if u.Path == "" {
		u.Path = "/"
	}

	return &Request{
		URL: u,
	}, nil
}

func Do(req *Request) ([]byte, error) {
	var err error
	var TlsTimeout time.Duration = time.Duration(15) * time.Second

	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, err
	}

	host := parsedUrl.Hostname()
	if host == "" {
		return nil, fmt.Errorf("incomplete url: %s", req.URL.String())
	}

	port := parsedUrl.Port()
	if parsedUrl.Port() == "" {
		port = "1965"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	conf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	cnx, err := tls.DialWithDialer(&net.Dialer{Timeout: TlsTimeout}, "tcp", addr, conf)
	if err != nil {
		log.Printf("failed to connect to %s: %v", parsedUrl.Host, err)
		return nil, err
	}
	defer cnx.Close()

	// Begin TOFU

	// End TOFU

	query := fmt.Sprintf("%s\r\n", req.URL.String())

	_, err = cnx.Write([]byte(query))
	if err != nil {
		return nil, err
	}

	result, err := io.ReadAll(cnx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FetchCapsuleGemini(request *Request) (Capsule, error) {
	capsule := Capsule{0, "", nil}

	result, err := Do(request)
	if err != nil {
		return capsule, err
	}

	resp := strings.SplitN(string(result), "\r\n", 2)
	if len(resp) < 2 {
		return capsule, fmt.Errorf("received an invalid response")
	}

	// Process the header
	header := strings.SplitN(resp[0], " ", 2)
	if len(header[0]) != 2 {
		header = strings.SplitN(resp[0], "\t", 2) // Do we need tab separation?
		if len(header[0]) != 2 {
			return capsule, fmt.Errorf("invalid response format")
		}
	}

	status, err := strconv.Atoi(string(header[0][0]))
	if err != nil {
		return Capsule{}, fmt.Errorf("invalid status code received: %v", err)
	}

	capsule.Status = status
	capsule.Information = header[1]
	capsule.Body = []byte(resp[1])

	return capsule, nil
}

func ParseGemText(body []byte, currentUrl string, active int) (string, []map[int]string, error) {
	documentStyle := lipgloss.NewStyle()
	//	Background(lipgloss.Color("#7D56F4"))

	// textStyle := lipgloss.NewStyle().
	// 	Inherit(documentStyle).
	// 	Width(constants.WindowWidth).
	// 	Foreground(lipgloss.Color("#FAFAFA"))

	typeStyle := lipgloss.NewStyle().
		Inherit(documentStyle)

	textStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Foreground(lipgloss.Color("#AEAEAE"))

	headingStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	linkStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	activeLinkStyle := lipgloss.NewStyle().
		Inherit(typeStyle).
		Bold(true).
		Foreground(lipgloss.Color("#CC56F4"))

	wrapped := WrapContent(string(body), 100)
	lines := strings.Split(wrapped, "\n")

	var link, text string
	var links []map[int]string

	spacer := "      "
	outputIndex := 0
	preformat := false

	for i, line := range lines {
		lines[i] = strings.Trim(line, "\r\n")

		if len(line) > 0 {
			if len(line) >= 3 && line[:3] == "```" && !preformat {
				preformat = true
				continue
			} else if len(line) >= 3 && line[:3] == "```" && preformat {
				preformat = false
				continue
			}

			if !preformat {
				if line[0] == '#' {
					line = headingStyle.Render(line)
				} else if len(line) > 3 && line[:2] == "=>" {
					subLn := strings.Trim(line[2:], "\r\n\t \a")
					split := strings.IndexAny(subLn, " \t")

					if split < 0 || len(subLn)-1 <= split {
						link = subLn
						text = subLn
					} else {
						link = strings.Trim(subLn[:split], "\r\n\t \a")
						text = strings.Trim(subLn[split:], "\r\n\t \a")
					}
					line = linkStyle.Render(text)
					if i == active || active == 0 {
						line = activeLinkStyle.Render(text)
						active = -1
					}

					if !strings.Contains(link, "://") {
						base, err := url.Parse(currentUrl)
						if err != nil {
							continue
						}

						href, err := url.Parse(link)
						if err != nil {
							continue
						}

						link = base.ResolveReference(href).String()
					}

					links = append(links, map[int]string{i: link})
				} else {
					line = textStyle.Render(line)
				}
			} else {
				line = textStyle.Render(line)
			}
		}

		lines[outputIndex] = lipgloss.JoinHorizontal(lipgloss.Top, typeStyle.Render(spacer), line)

		outputIndex++
	}

	return strings.Join(lines[:outputIndex], "\n"), links, nil
}

func WrapContent(raw string, width int) string {
	width = min(width, 100)
	counter := 0
	var content strings.Builder
	content.Grow(len(raw))

	spacer := ""

	for i, ch := range raw {
		if ch == '\n' || ch == '\u0085' || ch == '\u2028' || ch == '\u2029' {
			content.WriteRune('\n')
			counter = 0
		} else if ch == '\r' || ch == '\v' || ch == '\b' || ch == '\f' || ch == '\a' {
			// Get rid of control characters we don't want
			continue
		} else if ch == '\t' {
			if counter < width {
				content.WriteRune(ch)
				// counter += 4
			} else {
				content.WriteRune('\n')
				counter = 0
			}
		} else if ch == ' ' {
			// Peek ahead if the next space is going to overflow the width
			for j, next := range raw[i+1:] {
				if next == ' ' {
					if (counter + j) >= width {
						content.WriteRune('\n')
						counter = 0
						content.WriteString(spacer)
						counter += len(spacer)
					} else {
						content.WriteRune(' ')
						counter++
					}
					break
				}
			}
		} else {
			if counter <= width {
				content.WriteRune(ch)
				counter++
			} else {
				content.WriteRune('\n')
				counter = 0
				content.WriteString(spacer)
				counter += len(spacer)
				content.WriteRune(ch)
				counter++
			}
		}
	}

	return content.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
