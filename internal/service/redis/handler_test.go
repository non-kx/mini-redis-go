package redis

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

func getRedisRequestPayload(t *testing.T, cmd uint8, key string, val tlv.TypeLengthValue) payload.RequestPayload {
	bod := payload.RedisRequestBody{
		Key:   key,
		Value: val,
	}
	rawbod, err := bod.ToTLV()

	assert.Nil(t, err)

	return payload.RequestPayload{
		Cmd:  cmd,
		Body: rawbod,
	}
}

func TestHandleGetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	plkey := "\x00\x00\x00\x00\x00\x00\x00\x00test_key"
	val := tlv.String("test_val")
	rawval, _ := val.ToTLV()
	pl := getRedisRequestPayload(t, payload.GetCmd, key, []byte{})

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	ctx.EXPECT().GetRedis(plkey).Times(1).Return(rawval)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawval)

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandleSetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	plkey := "\x00\x00\x00\x00\x00\x00\x00\x00test_key"
	val := tlv.String("test_val")
	rawval, _ := val.ToTLV()
	rawok := []byte{2, 0, 0, 0, 2, 79, 75}
	pl := getRedisRequestPayload(t, payload.SetCmd, key, rawval)

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	ctx.EXPECT().SetRedis(plkey, rawval).Times(1)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawok)

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}
