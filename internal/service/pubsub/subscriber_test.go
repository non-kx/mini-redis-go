package pubsub

import (
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
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
