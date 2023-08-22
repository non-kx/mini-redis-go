package model

import (
	"net"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewStringTopic(t *testing.T) {
	topic := NewTopic[*tlv.String]("test_topic")

	assert.NotNil(t, topic)
	assert.NotNil(t, topic.ConnDb)
	assert.Equal(t, "test_topic", topic.Name)
}

func TestAddConnToTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	addr := mocknet.NewMockAddr(ctrl)

	name := "test_topic"
	topic := &Topic[*tlv.String]{
		Name:   name,
		ConnDb: db.NewKVStore[net.Conn](nil),
	}

	assert.Equal(t, 0, len(topic.ConnDb.Storage))

	conn.EXPECT().RemoteAddr().Times(1).Return(addr)
	addr.EXPECT().String().Times(1)
	topic.AddConn(conn)

	assert.Equal(t, 1, len(topic.ConnDb.Storage))
}

func TestRemoveConnFromTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	addr := mocknet.NewMockAddr(ctrl)

	addr.EXPECT().String().Times(1)

	name := "test_topic"
	db := &db.KVStore[net.Conn]{
		Storage: map[string]net.Conn{addr.String(): conn},
	}
	topic := &Topic[*tlv.String]{
		Name:   name,
		ConnDb: db,
	}

	assert.Equal(t, 1, len(topic.ConnDb.Storage))

	conn.EXPECT().RemoteAddr().Times(1).Return(addr)
	addr.EXPECT().String().Times(1)
	topic.RemoveConn(conn)

	assert.Equal(t, 0, len(topic.ConnDb.Storage))
}

func TestTopicShouldInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocknet.NewMockConn(ctrl)
	addr := mocknet.NewMockAddr(ctrl)

	addr.EXPECT().String().Times(1)

	name := "test_topic"
	db := &db.KVStore[net.Conn]{
		Storage: map[string]net.Conn{addr.String(): conn},
	}
	topic := &Topic[*tlv.String]{
		Name:   name,
		ConnDb: db,
	}

	didinit := topic.DidInit()

	assert.True(t, didinit)
}

func TestTopicShouldNotInit(t *testing.T) {
	name := "test_topic"
	topic := &Topic[*tlv.String]{
		Name: name,
	}

	didinit := topic.DidInit()

	assert.False(t, didinit)
}
