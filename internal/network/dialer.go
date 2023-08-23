package network

import (
	"crypto/tls"
	"net"
)

type Dialer interface {
	Dial(network string, url string) (net.Conn, error)
	SecureDial(network string, url string, conf *tls.Config) (net.Conn, error)
}

type TcpDialer struct {
}

func (d *TcpDialer) Dial(network string, url string) (net.Conn, error) {
	return net.Dial(network, url)
}

func (d *TcpDialer) SecureDial(network string, url string, conf *tls.Config) (net.Conn, error) {
	return tls.Dial(network, url, conf)
}
