package redis

import (
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSendGetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	key := "test_key"

	reqbod := payload.RedisRequestBody{
		Key: key,
	}
	rawreqbod, _ := reqbod.ToTLV()
	test.ExpectWriteRequestToConn(t, conn, payload.GetCmd, rawreqbod)

	resbod := tlv.String("value")
	rawresbod, _ := resbod.ToTLV()
	test.ExpectReadResponseFromConn(t, conn, tlv.StringType, rawresbod)

	val, err := SendGetRequest(conn, key)

	assert.Nil(t, err)
	assert.Equal(t, "value", val.String())
}

func TestSendSetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	key := "test_key"
	val := tlv.String("value")
	rawval, _ := val.ToTLV()

	reqbod := payload.RedisRequestBody{
		Key:   key,
		Value: rawval,
	}
	rawreqbod, _ := reqbod.ToTLV()
	test.ExpectWriteRequestToConn(t, conn, payload.SetCmd, rawreqbod)

	resbod := tlv.String("OK")
	rawresbod, _ := resbod.ToTLV()
	test.ExpectReadResponseFromConn(t, conn, tlv.StringType, rawresbod)

	resp, err := SendSetRequest(conn, key, rawval)

	assert.Nil(t, err)
	assert.Equal(t, "OK", resp)
}
