package client

import (
	"context"
	"fmt"
	"net"

	"github.com/shripadmhetre/distrib-cache-golang/service"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &service.CommandGet{
		Key: key,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := service.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == service.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key (%s)", key)
	}
	if resp.Status != service.StatusOK {
		return nil, fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &service.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := service.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}
	if resp.Status != service.StatusOK {
		return fmt.Errorf("server responsed with non OK status [%s]", resp.Status)
	}

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
