package network

import (
	"errors"
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mocknetwork "bitbucket.org/non-pn/mini-redis-go/internal/mock/network"
	mockservice "bitbucket.org/non-pn/mini-redis-go/internal/mock/service"
	"bitbucket.org/non-pn/mini-redis-go/internal/service"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pubsub"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewClient(t *testing.T) {
	client := NewClient("tcp", "localhost", ":6337", "", "")

	assert.NotNil(t, client)
	assert.Equal(t, "tcp", client.Network)
	assert.Equal(t, "localhost", client.Host)
	assert.Nil(t, client.Connection)
	assert.NotNil(t, client.Service)
}

func TestClientConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	transport := mocknetwork.NewMockTransporter(ctrl)
	conn := mocknet.NewMockConn(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Transport:  transport,
		Connection: nil,
		Service:    &service.Service{},
	}

	transport.EXPECT().EstablishConnection().Times(1).Return(conn, nil)

	err := client.Connect()

	assert.Nil(t, err)
	assert.Equal(t, conn, client.Connection)
}

func TestClientClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhots",
		Connection: conn,
		Service:    &service.Service{},
	}

	conn.EXPECT().Close().Times(1)

	err := client.Close()

	assert.Nil(t, err)
	assert.Nil(t, client.Connection)
}

func TestClientPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	pong := tlv.String("PONG")
	service.EXPECT().SendPingRequest(conn, nil).Times(1).Return(&pong, nil)

	resp, err := client.Ping(nil)

	assert.Nil(t, err)
	assert.Equal(t, pong.String(), resp)
}

func TestClientGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	key := "test_key"
	val := tlv.String("test_val")
	service.EXPECT().SendGetRequest(conn, key).Times(1).Return(&val, nil)

	resp, err := client.Get(key)

	assert.Nil(t, err)
	assert.Equal(t, val.String(), resp.String())
}

func TestClientGetNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	key := "test_key"
	service.EXPECT().SendGetRequest(conn, key).Times(1).Return(nil, nil)

	resp, err := client.Get(key)

	assert.Nil(t, err)
	assert.Equal(t, "", resp.String())
}

func TestClientGetReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	key := "test_key"
	mockerr := errors.New("Some send get request error")
	service.EXPECT().SendGetRequest(conn, key).Times(1).Return(nil, mockerr)

	resp, err := client.Get(key)

	assert.NotNil(t, err)
	assert.Equal(t, nil, resp)
}

func TestClientSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	key := "test_key"
	val := tlv.String("test_val")
	raw := []byte{2, 0, 0, 0, 8, 116, 101, 115, 116, 95, 118, 97, 108}
	service.EXPECT().SendSetRequest(conn, key, raw).Times(1).Return("OK", nil)

	resp, err := client.Set(key, &val)

	assert.Nil(t, err)
	assert.Equal(t, "OK", resp)
}

func TestClientSetReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	key := "test_key"
	val := tlv.String("test_val")
	mockerr := errors.New("Some send set request error")
	service.EXPECT().SendSetRequest(conn, gomock.Any(), gomock.Any()).Times(1).Return("", mockerr)

	resp, err := client.Set(key, &val)

	assert.NotNil(t, err)
	assert.Equal(t, "", resp)
}

func TestClientSub(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	sub := &pubsub.Subscriber{
		Conn:         conn,
		IsSubscribed: false,
	}

	topic := "test_topic"
	service.EXPECT().SendSubRequest(conn, topic).Times(1).Return(sub, nil)

	sub, err := client.Sub(topic)

	assert.Nil(t, err)
	assert.Equal(t, sub, sub)
}

func TestClientSubReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	topic := "test_topic"
	mockerr := errors.New("Some send sub request error")
	service.EXPECT().SendSubRequest(conn, topic).Times(1).Return(nil, mockerr)

	sub, err := client.Sub(topic)

	assert.Equal(t, mockerr, err)
	assert.Nil(t, sub)
}

func TestClientPub(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	topic := "test_topic"
	val := "test_val"
	raw := []byte{2, 0, 0, 0, 8, 116, 101, 115, 116, 95, 118, 97, 108}
	service.EXPECT().SendPubRequest(conn, topic, raw).Times(1).Return("OK", nil)

	resp, err := client.Pub(topic, val)

	assert.Nil(t, err)
	assert.Equal(t, "OK", resp)
}

func TestClientPubReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	service := mockservice.NewMockIService(ctrl)

	client := &Client{
		Network:    "tcp",
		Host:       "localhost",
		Connection: conn,
		Service:    service,
	}

	topic := "test_topic"
	val := "test_val"
	raw := []byte{2, 0, 0, 0, 8, 116, 101, 115, 116, 95, 118, 97, 108}
	mockerr := errors.New("Some send pub request error")
	service.EXPECT().SendPubRequest(conn, topic, raw).Times(1).Return("", mockerr)

	resp, err := client.Pub(topic, val)

	assert.Equal(t, mockerr, err)
	assert.Equal(t, "", resp)
}
