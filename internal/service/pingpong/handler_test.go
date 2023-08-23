package pingpong

import (
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mockpayload "bitbucket.org/non-pn/mini-redis-go/internal/mock/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlePingRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	pingstr := tlv.String("PING")
	pingbod, _ := pingstr.ToTLV()
	pl := &payload.RequestPayload{
		Cmd:  payload.PingCmd,
		Body: pingbod,
	}
	ctx.EXPECT().GetPayload().Times(1).Return(pl)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteStringResponseToConn(t, conn, "PONG")

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandleCustomPingRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	pingstr := tlv.String("Custom")
	pingbod, _ := pingstr.ToTLV()
	pl := &payload.RequestPayload{
		Cmd:  payload.PingCmd,
		Body: pingbod,
	}
	ctx.EXPECT().GetPayload().Times(1).Return(pl)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteStringResponseToConn(t, conn, "Custom")

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}
