package pubsub

import (
	"errors"
	"io"
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/test"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getTestSubscriber(t *testing.T) (*Subscriber, *mocknet.MockConn) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	return &Subscriber{
		Conn:         conn,
		IsSubscribed: false,
	}, conn
}

func TestNewSubscriber(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)

	sub := NewSubscriber(conn)

	assert.NotNil(t, sub)
	assert.NotNil(t, sub.Conn)
	assert.False(t, sub.IsSubscribed)
}

func TestNextMesssageIsSubscribe(t *testing.T) {
	sub, _ := getTestSubscriber(t)
	sub.IsSubscribed = true
	msg, err := sub.NextMessage()

	assert.Nil(t, msg)
	assert.Nil(t, err)
}

func TestNextMesssageIsNotSubscribe(t *testing.T) {
	sub, conn := getTestSubscriber(t)

	typ := []byte{0x3}
	blen := []byte{0x0, 0x0, 0x0, 0xe}
	val := []byte{0xc, 0x2, 0x0, 0x0, 0x0, 0x8, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x73, 0x67}

	conn.EXPECT().Read(gomock.Any()).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, val).Return(len(val), nil)
	msg, err := sub.NextMessage()

	assert.Nil(t, err)
	assert.Equal(t, "test_msg", msg.String())
}

func TestNextMesssageEOF(t *testing.T) {
	sub, conn := getTestSubscriber(t)

	conn.EXPECT().Read(gomock.Any()).Return(0, io.EOF)

	msg, err := sub.NextMessage()

	assert.Nil(t, msg)
	assert.Equal(t, io.EOF, err)
}

func TestNextMesssageIsNotMsgType(t *testing.T) {
	sub, conn := getTestSubscriber(t)

	typ := []byte{0x3}
	blen := []byte{0x0, 0x0, 0x0, 0xe}
	val := []byte{0x1, 0x2, 0x0, 0x0, 0x0, 0x8, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x73, 0x67} // First byte is not 0xc(Msg type)

	conn.EXPECT().Read(gomock.Any()).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, val).Return(len(val), nil)
	msg, err := sub.NextMessage()

	assert.Equal(t, errors.New("Error not message type"), err)
	assert.Nil(t, msg)
}

func TestSubscribeToTopic(t *testing.T) {
	sub, conn := getTestSubscriber(t)

	msg := tlv.String("test_msg")
	raw, _ := msg.ToTLV()
	test.ExpectReadResponseFromConn(t, conn, tlv.MsgType, raw)

	conn.EXPECT().Read(gomock.Any()).Return(0, io.EOF)

	sub.Subscribe(func(s string) {
		assert.Equal(t, s, msg.String())
	})
}
