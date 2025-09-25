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

func ParseGemText(body []byte, currentUrl string) string {
	lines := strings.Split(string(body), "\n")

	// spacer := "      "
	spacer := ""

	outputIndex := 0
	for i, line := range lines {
		lines[i] = strings.Trim(line, "\r\n")

		var leader, tail string = "", ""
		if len(line) > 0 && line[0] == '#' {
			leader = "\033[1m"
			tail = "\033[0m"
		}
		lines[outputIndex] = fmt.Sprintf("%s%s%s%s", spacer, leader, line, tail)

		outputIndex++
	}

	foo := strings.Join(lines[:outputIndex], "\n")

	return WrapContent(foo, 80)
}

func WrapContent(raw string, width int) string {
	width = min(width, 80)
	counter := 0
	var content strings.Builder
	content.Grow(len(raw))

	//spacer := "      "
	log.Println(width)

	for _, ch := range raw {
		if ch == '\n' || ch == '\u0085' || ch == '\u2028' || ch == '\u2029' {
			content.WriteRune('\n')
			counter = 0
		} else if ch == '\r' || ch == '\v' || ch == '\b' || ch == '\f' || ch == '\a' {
			// Get rid of control characters we dont want
			continue
		} else if ch == '\t' {
			if counter < width {
				content.WriteString("")
				// counter += 4
			} else {
				content.WriteRune('\n')
				counter = 0
			}
		} else {
			if counter <= width {
				content.WriteRune(ch)
				counter++
			} else {
				content.WriteRune('\n')
				counter = 0
				// content.WriteString(spacer)
				// counter += len(spacer)
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
