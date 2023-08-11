package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type IClient interface {
	Connect(network string, host string) error
	Close() error
	Send(data []byte) ([]byte, error)
}

type Client struct {
	Network    string
	Host       string
	Connection *net.Conn
}

func (c *Client) Connect() error {
	conn, err := net.Dial(c.Network, c.Host)
	if err != nil {
		return err
	}

	c.Connection = &conn

	return nil
}

func (c *Client) Close() error {
	err := (*c.Connection).Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Send(data []byte) ([]byte, error) {
	fmt.Fprintf(*c.Connection, string(data)+"\n")
	// (*c.Connection).Write(data)

	resp, err := bufio.NewReader(*c.Connection).ReadString('\n')
	if err != nil {
		return nil, err
	}

	return []byte(strings.TrimSuffix(resp, "\n")), nil
}

func NewClient(network string, host string) *Client {
	return &Client{
		Network:    network,
		Host:       host,
		Connection: nil,
	}
}
