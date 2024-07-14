package client

import (
	"bytes"
	"context"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(address string) *Client {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		addr: address,
		conn: conn,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)

	wr.WriteArray(
		[]resp.Value{resp.StringValue("SET"),
			resp.StringValue(key),
			resp.StringValue(val)},
	)

	_, err = conn.Write(buf.Bytes())
	// io.Copy(c.conn, buf)
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)

	wr.WriteArray(
		[]resp.Value{resp.StringValue("GET"),
			resp.StringValue(key),
		})

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return "", err

	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	// io.Copy(c.conn, buf)
	return string(b[:n]), err
}
