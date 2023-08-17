package network

import "net"

type Transporter interface {
	net.Listener
	net.Dialer
}

// func GetTransport(cert string, key string) (*Transporter, error) {

// }
