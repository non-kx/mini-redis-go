package pubsub

import (
	"net"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	mockpayload "bitbucket.org/non-pn/mini-redis-go/internal/mock/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getPubsubRequestPayload(t *testing.T, cmd uint8, topic string, val tlv.TypeLengthValue) payload.RequestPayload {
	bod := payload.PubsubRequestBody{
		Topic: topic,
		Value: val,
	}
	rawbod, err := bod.ToTLV()

	assert.Nil(t, err)

	return payload.RequestPayload{
		Cmd:  cmd,
		Body: rawbod,
	}
}

func TestHandleSubRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	addr := mocknet.NewMockAddr(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	topicname := "test_topic"
	pltopic := "\x00\x00\x00\x00\x00\x00test_topic"
	pl := getPubsubRequestPayload(t, payload.SubCmd, topicname, []byte{})

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	conn.EXPECT().RemoteAddr().Times(1).Return(addr)
	ctx.EXPECT().GetPubsub(pltopic).Times(1).Return(nil)
	addr.EXPECT().String().Times(1).Return("localhost")
	ctx.EXPECT().SetPubsub(pltopic, gomock.Any()).Times(1)
	ctx.EXPECT().GetConn().AnyTimes().Return(conn)

	rawok := []byte{2, 0, 0, 0, 2, 79, 75}
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawok)

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}

func TestHandlePubRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	subconn := mocknet.NewMockConn(ctrl)
	subaddr := mocknet.NewMockAddr(ctrl)
	ctx := mockpayload.NewMockIRequestContext(ctrl)

	topicname := "test_topic"
	pltopic := "\x00\x00\x00\x00\x00\x00test_topic"
	msg := []byte{2, 0, 0, 0, 8, 116, 101, 115, 116, 95, 109, 115, 103}
	pl := getPubsubRequestPayload(t, payload.PubCmd, topicname, msg)

	topic := &model.Topic[*tlv.String]{
		Name: pltopic,
		ConnDb: &db.KVStore[net.Conn]{
			Storage: map[string]net.Conn{"localhost": subconn},
		},
	}

	ctx.EXPECT().GetPayload().Times(1).Return(&pl)
	ctx.EXPECT().GetPubsub(pltopic).Times(1).Return(topic)
	ctx.EXPECT().GetConn().AnyTimes().Return(conn)

	subconn.EXPECT().RemoteAddr().Times(1).Return(subaddr)
	subaddr.EXPECT().String().Times(1).Return("localhost")
	test.ExpectWriteResponseToConn(t, subconn, tlv.MsgType, msg)

	rawok := []byte{2, 0, 0, 0, 2, 79, 75}
	test.ExpectWriteResponseToConn(t, conn, tlv.StringType, rawok)

	err := HandleRequest(ctx)

	assert.Nil(t, err)
}
