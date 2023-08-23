package network

import (
	"errors"
	"io"
	"net"
	"path/filepath"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mockservice "bitbucket.org/non-pn/mini-redis-go/internal/mock/service"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func expectReadRequestFromConn(t *testing.T, conn *mocknet.MockConn, cmd uint8, body tlv.TypeLengthValue) {
	req := payload.RequestPayload{
		Cmd:  cmd,
		Body: body,
	}
	raw, err := req.ToTLV()
	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Read(gomock.Any()).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, val).Return(len(val), nil)
}

func TestNewServer(t *testing.T) {
	serv, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "", "")
	defer serv.Stop()

	assert.Nil(t, err)
	assert.NotNil(t, serv)
	assert.NotNil(t, serv.RedisDb)
	assert.NotNil(t, serv.PubsubDb)
	assert.Equal(t, 0, len(serv.Connections))
}

func TestNewSecureServerWithInvalidCert(t *testing.T) {
	_, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "invalid", "invalid")

	assert.NotNil(t, err)
}

func TestNewSecureServerWithValidCert(t *testing.T) {
	cert, err := filepath.Abs("../../test/cert/server/server.crt")
	key, err := filepath.Abs("../../test/cert/server/server.key")

	assert.Nil(t, err)

	serv, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, cert, key)
	defer serv.Stop()

	assert.Nil(t, err)
	assert.NotNil(t, serv)
	assert.NotNil(t, serv.RedisDb)
	assert.NotNil(t, serv.PubsubDb)
	assert.Equal(t, 0, len(serv.Connections))
}

func TestStopServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	listener := mocknet.NewMockListener(ctrl)
	conn := mocknet.NewMockConn(ctrl)

	server := &Server{
		Port:        ":6337",
		Listener:    listener,
		Connections: []net.Conn{conn},
		Service:     &service.Service{},
		RedisDb:     &db.KVStore[[]byte]{},
		PubsubDb:    &db.KVStore[*model.Topic[*tlv.String]]{},
	}

	conn.EXPECT().Close().Times(1)
	listener.EXPECT().Close().Times(1)

	err := server.Stop()

	assert.Nil(t, err)
}

func TestHandleConnectionHandleRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	listener := mocknet.NewMockListener(ctrl)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	server := &Server{
		Port:        ":6337",
		Listener:    listener,
		Connections: []net.Conn{},
		Service:     service,
		RedisDb:     &db.KVStore[[]byte]{},
		PubsubDb:    &db.KVStore[*model.Topic[*tlv.String]]{},
	}

	expectReadRequestFromConn(t, conn, payload.GetCmd, []byte{})

	mockerr := errors.New("Some handle request error")
	service.EXPECT().HandleRequest(gomock.Any()).AnyTimes().Return(mockerr)

	err := server.HandleConnection(conn)

	assert.Equal(t, mockerr, err)
}

func TestHandleConnectionEOF(t *testing.T) {
	ctrl := gomock.NewController(t)
	listener := mocknet.NewMockListener(ctrl)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	server := &Server{
		Port:        ":6337",
		Listener:    listener,
		Connections: []net.Conn{},
		Service:     service,
		RedisDb:     &db.KVStore[[]byte]{},
		PubsubDb:    &db.KVStore[*model.Topic[*tlv.String]]{},
	}
	service.EXPECT().HandleRequest(conn).AnyTimes()
	conn.EXPECT().Read(gomock.Any()).Return(0, io.EOF)

	err := server.HandleConnection(conn)

	assert.Equal(t, io.EOF, err)
}
