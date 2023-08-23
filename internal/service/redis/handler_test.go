package redis

import (
	"errors"
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mockpayload "bitbucket.org/non-pn/mini-redis-go/internal/mock/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getRedisRequestBody(key string, val tlv.TypeLengthValue) payload.RedisRequestBody {
	return payload.RedisRequestBody{
		Key:   key,
		Value: val,
	}
}

func getRedisRequestPayload(t *testing.T, cmd uint8, key string, val tlv.TypeLengthValue) payload.RequestPayload {
	bod := getRedisRequestBody(key, val)
	rawbod, err := bod.ToTLV()

	assert.Nil(t, err)

	return payload.RequestPayload{
		Cmd:  cmd,
		Body: rawbod,
	}
}

func TestHandleRedisRequestGetCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	pl := getRedisRequestPayload(t, payload.GetCmd, key, []byte{})

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	ctx.EXPECT().GetRedis(gomock.Any()).Times(1)
	ctx.EXPECT().GetConn().AnyTimes().Return(conn)
	conn.EXPECT().Write(gomock.Any()).AnyTimes()

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandleRedisRequestSetCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	pl := getRedisRequestPayload(t, payload.SetCmd, key, []byte{})

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	ctx.EXPECT().SetRedis(gomock.Any(), gomock.Any()).Times(1)
	ctx.EXPECT().GetConn().AnyTimes().Return(conn)
	conn.EXPECT().Write(gomock.Any()).AnyTimes()

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandleRedisRequestNotMatchCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	defaultcmd := uint8(0)
	pl := getRedisRequestPayload(t, defaultcmd, key, []byte{})

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandleRedisRequestInvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	pl := getRedisRequestPayload(t, payload.SetCmd, key, []byte{})
	invalidbod := []byte{7, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0, 0, 116, 101, 115, 116, 95, 107, 101, 121} // Tlv invalid data type
	pl.Body = invalidbod

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)

	err := HandleRequest(ctx)

	assert.Equal(t, errors.New("Invalid redis request"), err)
}

func TestHandleGetRequestStringType(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.String("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, []byte{})

	ctx.EXPECT().GetRedis(bod.Key).Times(1).Return(rawval)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawval)

	err := handleGetRequest(ctx, &bod)

	assert.Nil(t, err)
}

func TestHandleGetRequestBinaryType(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.Binary("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, []byte{})

	ctx.EXPECT().GetRedis(bod.Key).Times(1).Return(rawval)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.BinaryType, rawval)

	err := handleGetRequest(ctx, &bod)

	assert.Nil(t, err)
}

func TestHandleGetRequestReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.Binary("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, []byte{})

	ctx.EXPECT().GetRedis(bod.Key).Times(1).Return(rawval)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	conn.EXPECT().Write(gomock.Any()).AnyTimes().Return(0, errors.New("Some write error"))

	err := handleGetRequest(ctx, &bod)

	assert.NotNil(t, err)
}

func TestHandleSetRequestStringType(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.String("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, rawval)
	rawok := []byte{2, 0, 0, 0, 2, 79, 75}

	ctx.EXPECT().SetRedis(bod.Key, rawval).Times(1)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawok)

	err := handleSetRequest(ctx, &bod)

	assert.Nil(t, err)
}

func TestHandleSetRequestBinaryType(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.Binary("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, rawval)
	rawok := []byte{2, 0, 0, 0, 2, 79, 75}

	ctx.EXPECT().SetRedis(bod.Key, rawval).Times(1)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawok)

	err := handleSetRequest(ctx, &bod)

	assert.Nil(t, err)
}

func TestHandleSetRequestWriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	key := "test_key"
	val := tlv.Binary("test_val")
	rawval, _ := val.ToTLV()
	bod := getRedisRequestBody(key, rawval)

	ctx.EXPECT().SetRedis(bod.Key, rawval).Times(1)
	ctx.EXPECT().GetConn().Times(1).Return(conn)
	conn.EXPECT().Write(gomock.Any()).AnyTimes().Return(0, errors.New("Some write error"))

	err := handleSetRequest(ctx, &bod)

	assert.NotNil(t, err)
}
