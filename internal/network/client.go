package network

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/service/pingpong"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pubsub"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
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

func (c *Client) Connect(cert string, key string) error {
	conn, err := EstablishConnection(c.Network, c.Host, cert, key)
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

func (c *Client) Ping(msg *string) (string, error) {
	resp, err := pingpong.SendPingRequest(c.Connection, msg)
	if err != nil {
		return "", err
	}

	return string(*resp), nil
}

func (c *Client) Get(k string) (tlv.TLVCompatible, error) {
	resp, err := redis.SendGetRequest(c.Connection, k)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return new(tlv.String), nil
	}

	return resp, nil
}

func (c *Client) Set(k string, v tlv.TLVCompatible) (string, error) {
	raw, err := v.ToTLV()
	if err != nil {
		return "", err
	}

	resp, err := redis.SendSetRequest(c.Connection, k, raw)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *Client) Sub(topic string) (*pubsub.Subscriber, error) {
	sub, err := pubsub.SendSubRequest(c.Connection, topic)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (c *Client) Pub(topic string, msg string) (string, error) {
	s := tlv.String(msg)
	raw, err := s.ToTLV()
	if err != nil {
		return "", err
	}

	resp, err := pubsub.SendPubRequest(c.Connection, topic, raw)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func NewClient(network string, host string) *Client {
	return &Client{
		Network:    network,
		Host:       host,
		Connection: nil,
	}
}
