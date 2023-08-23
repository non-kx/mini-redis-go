package network

import (
	"path/filepath"
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mocknetwork "bitbucket.org/non-pn/mini-redis-go/internal/mock/network"
	"bitbucket.org/non-pn/mini-redis-go/internal/network/ssl"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEstablishConnection(t *testing.T) {
	// cert, err := filepath.Abs("../../test/cert/client/client.crt")
	// key, err := filepath.Abs("../../test/cert/client/client.key")
	ctrl := gomock.NewController(t)
	dialer := mocknetwork.NewMockDialer(ctrl)
	conn := mocknet.NewMockConn(ctrl)

	transport := &TcpTransport{
		Network:  "tcp",
		Host:     "localhost",
		Port:     ":6377",
		CertPath: "",
		KeyPath:  "",
		Dialer:   dialer,
	}

	dialer.EXPECT().Dial(transport.Network, transport.Host+transport.Port).Times(1).Return(conn, nil)

	c, err := transport.EstablishConnection()

	assert.Nil(t, err)
	assert.NotNil(t, conn, c)
}

func TestEstablishSecureConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	dialer := mocknetwork.NewMockDialer(ctrl)
	conn := mocknet.NewMockConn(ctrl)

	cert, err := filepath.Abs("../../test/cert/client/client.crt")
	key, err := filepath.Abs("../../test/cert/client/client.key")

	transport := &TcpTransport{
		Network:  "tcp",
		Host:     "localhost",
		Port:     ":6377",
		CertPath: cert,
		KeyPath:  key,
		Dialer:   dialer,
	}

	conf, err := ssl.GetClientTlsConfig(cert, key)

	assert.Nil(t, err)

	dialer.EXPECT().SecureDial(transport.Network, transport.Host+transport.Port, conf).Times(1).Return(conn, nil)

	c, err := transport.EstablishConnection()

	assert.Nil(t, err)
	assert.NotNil(t, conn, c)
}
