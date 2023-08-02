package network

import (
	"bufio"
	"fmt"
	"net"
)

type IClient interface {
	Connect(network string, host string) error
	Send(msg string) ([]byte, error)
	Close() error
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

func (c *Client) Send(msg string) ([]byte, error) {
	fmt.Fprintf(*c.Connection, msg+"\n")

	resp, err := bufio.NewReader(*c.Connection).ReadString('\n')
	if err != nil {
		return nil, err
	}

	return []byte(resp), nil
}

func (c *Client) Close() error {
	err := (*c.Connection).Close()
	if err != nil {
		return err
	}

	return nil
}

func NewClient(network string, host string) *Client {
	return &Client{
		Network:    network,
		Host:       host,
		Connection: nil,
	}
}
