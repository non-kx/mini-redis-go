package pingpong

import (
	"errors"
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSendPingRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	msg := "PING"
	rawping := []byte{2, 0, 0, 0, 4, 80, 73, 78, 71}

	test.ExpectWriteRequestToConn(t, conn, payload.PingCmd, rawping)
	test.ExpectReadStringResponseFromConn(t, conn, "PONG")

	resp, err := SendPingRequest(conn, &msg)

	assert.Nil(t, err)
	assert.Equal(t, resp.String(), "PONG")
}

func TestSendPingRequestWithCustomMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	msg := "Custom"
	s := tlv.String(msg)
	reqbod, _ := s.ToTLV()

	test.ExpectWriteRequestToConn(t, conn, payload.PingCmd, reqbod)
	test.ExpectReadStringResponseFromConn(t, conn, msg)

	resp, err := SendPingRequest(conn, &msg)

	assert.Nil(t, err)
	assert.Equal(t, resp.String(), msg)
}

func TestSendPingRequestWithNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	rawping := []byte{2, 0, 0, 0, 4, 80, 73, 78, 71}

	test.ExpectWriteRequestToConn(t, conn, payload.PingCmd, rawping)
	test.ExpectReadStringResponseFromConn(t, conn, "PONG")

	resp, err := SendPingRequest(conn, nil)

	assert.Nil(t, err)
	assert.Equal(t, resp.String(), "PONG")
}

func TestSendPingRequestWithWriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	conn.EXPECT().Write(gomock.Any()).Return(0, errors.New("Some write error"))

	msg := "PING"
	_, err := SendPingRequest(conn, &msg)

	assert.NotNil(t, err)
}

func TestSendPingRequestWithReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	conn.EXPECT().Write(gomock.Any()).AnyTimes()
	conn.EXPECT().Read(gomock.Any()).Return(0, errors.New("Some read error"))

	msg := "PING"
	_, err := SendPingRequest(conn, &msg)

	assert.NotNil(t, err)
}
