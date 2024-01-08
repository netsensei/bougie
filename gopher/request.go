package gopher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
)

type Request struct {
	Host     string
	Port     string
	ItemType string
	Selector string
	Query    string
}

func New() *Request {
	return &Request{
		Port:     "70",
		ItemType: "1",
		Selector: "\r\n",
		Query:    "",
	}
}

func (r *Request) Do(ctx context.Context) ([]byte, error) {
	var d net.Dialer

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

	types := []string{"0", "1", "3", "7"}
	if !slices.Contains(types, r.ItemType) {
		return nil, errors.New("unsupported item type for gopher request")
	}

	return data, nil
}
