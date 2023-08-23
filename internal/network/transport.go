package network

import (
	"crypto/tls"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/network/ssl"
)

type Transporter interface {
	GetListener() (net.Listener, error)
	EstablishConnection() (net.Conn, error)
}

type TcpTransport struct {
	Network  string
	Host     string
	Port     string
	CertPath string
	KeyPath  string
	Dialer   Dialer
}

func (t *TcpTransport) GetListener() (net.Listener, error) {
	var (
		l             net.Listener
		cert, key     = t.CertPath, t.KeyPath
		network, port = t.Network, t.Port

		err error
	)
	if cert == "" || key == "" {
		l, err = net.Listen(network, port)
		if err != nil {
			return nil, err
		}

		return l, nil
	}

	log.Println("Use secure tcp")
	tlsconf, err := ssl.GetServerTlsConfig(cert, key)
	l, err = tls.Listen(network, port, tlsconf)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (t *TcpTransport) EstablishConnection() (net.Conn, error) {
	var (
		conn      net.Conn
		cert, key = t.CertPath, t.KeyPath
		network   = t.Network
		url       = t.Host + t.Port

		err error
	)
	if cert == "" || key == "" {
		conn, err = t.Dialer.Dial(network, url)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	log.Println("Use secure tcp")
	tlsconf, err := ssl.GetClientTlsConfig(cert, key)
	conn, err = t.Dialer.SecureDial(network, url, tlsconf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewTcpTransport(network string, host string, port string, cert string, key string) Transporter {
	return &TcpTransport{
		Network:  network,
		Host:     host,
		Port:     port,
		CertPath: cert,
		KeyPath:  key,
		Dialer:   &TcpDialer{},
	}
}
