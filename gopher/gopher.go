package gopher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"slices"
	"strings"
)

type Request struct {
	Host     string
	Port     string
	ItemType string
	Selector string
	Query    string
}

type Response struct {
	Body     string
	ItemType string
	Request  *Request
}

func New(u string) *Request {
	p, _ := url.Parse(u)
	r := &Request{
		Host:     p.Hostname(),
		Port:     "70",
		ItemType: "1",
		Selector: "\r\n",
		Query:    "",
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
		r.Selector = "/" + strings.Join(parts[2:], "/") + "\r\n"
	}

	return r
}

func (r *Request) Do(ctx context.Context) (*Response, error) {
	var d net.Dialer

	types := []string{"0", "1", "3", "7"}
	if !slices.Contains(types, r.ItemType) {
		return nil, errors.New("unsupported item type for gopher request")
	}

	cnx, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", r.Host, r.Port))
	if err != nil {
		return nil, err
	}
	defer cnx.Close()

	fmt.Fprintf(cnx, "%s", r.Selector)
	data, err := io.ReadAll(io.LimitReader(cnx, 1024*1024))
	if err != nil {
		return nil, err
	}

	return &Response{
		Body:     string(data),
		ItemType: r.ItemType,
		Request:  r,
	}, nil
}
