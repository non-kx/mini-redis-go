package pingpong

import (
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
	s := tlv.String(msg)
	reqbod, _ := s.ToTLV()

	test.ExpectWriteRequestToConn(t, conn, payload.PingCmd, reqbod)
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
