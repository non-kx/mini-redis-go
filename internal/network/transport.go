package network

import (
	"crypto/tls"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/network/ssl"
)

func GetListener(network string, port string, cert string, key string) (net.Listener, error) {
	var (
		l   net.Listener
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

func EstablishConnection(network string, url string, cert string, key string) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)
	if cert == "" || key == "" {
		conn, err = net.Dial(network, url)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	log.Println("Use secure tcp")
	tlsconf, err := ssl.GetClientTlsConfig(cert, key)
	conn, err = tls.Dial(network, url, tlsconf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
