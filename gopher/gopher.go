package gopher

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Host     string
	Port     string
	ItemType string
	Selector string
	Query    string
}

type Response struct {
	Body     []byte
	ItemType string
}

func New(u string) *Request {
	p, _ := url.Parse(u)

	r := &Request{
		Host:     p.Hostname(),
		Port:     "70",
		ItemType: "1",
		Selector: "\r\n",
	}

	if len(p.Query()) > 0 {
		r.Query = p.Query().Get("q")
	}

	if len(p.Port()) > 0 {
		r.Port = p.Port()
	}

	parts := strings.Split(p.Path, "/")
	if len(parts) > 1 {
		r.ItemType = parts[1]
		items := strings.Split(parts[len(parts)-1], "%09")
		if len(items) > 1 {
			r.Query = items[1]
		}
		parts[len(parts)-1] = strings.Replace(parts[len(parts)-1], "%09", "\t", 1)
		r.Selector = "/" + strings.Join(parts[2:], "/")
	}

	return r
}

func (r *Request) Do(ctx context.Context) (*Response, error) {
	d := net.Dialer{
		Timeout: 15 * time.Second,
	}

	cnx, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", r.Host, r.Port))
	if err != nil {
		return nil, err
	}
	defer cnx.Close()

	if r.Query != "" {
		fmt.Fprintf(cnx, "%s\t%s%s", r.Selector, r.Query, "\r\n")
	} else {
		fmt.Fprintf(cnx, "%s%s", r.Selector, "\r\n")
	}

	data, err := io.ReadAll(io.LimitReader(cnx, 1024*1024))
	if err != nil {
		return nil, err
	}

	return &Response{
		Body:     data,
		ItemType: r.ItemType,
	}, nil
}
