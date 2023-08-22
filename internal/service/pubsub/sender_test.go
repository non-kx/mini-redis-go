package pubsub

import (
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSendSubRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	topic := "test_topic"

	reqbod := payload.PubsubRequestBody{
		Topic: topic,
	}
	rawreqbod, _ := reqbod.ToTLV()
	test.ExpectWriteRequestToConn(t, conn, payload.SubCmd, rawreqbod)

	resbod := tlv.String("OK")
	rawresbod, _ := resbod.ToTLV()
	test.ExpectReadResponseFromConn(t, conn, tlv.StringType, rawresbod)

	sub, err := SendSubRequest(conn, topic)

	assert.Nil(t, err)
	assert.NotNil(t, sub)
}

func TestSendPubRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	topic := "test_topic"
	msg := tlv.String("test_msg")
	rawmsg, _ := msg.ToTLV()

	reqbod := payload.PubsubRequestBody{
		Topic: topic,
		Value: rawmsg,
	}
	rawreqbod, _ := reqbod.ToTLV()
	test.ExpectWriteRequestToConn(t, conn, payload.PubCmd, rawreqbod)

	resbod := tlv.String("OK")
	rawresbod, _ := resbod.ToTLV()
	test.ExpectReadResponseFromConn(t, conn, tlv.StringType, rawresbod)

	resp, err := SendPubRequest(conn, topic, rawmsg)

	assert.Nil(t, err)
	assert.Equal(t, "OK", resp)
}
