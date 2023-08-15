package payload

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

const (
	TOPIC_CHAN_CAP = 10
)

type TopicChan[T tlv.TLVCompatible] chan T

func NewTopicChan[T tlv.TLVCompatible](topic string) *TopicChan[T] {
	c := make(chan T, TOPIC_CHAN_CAP)
	tc := TopicChan[T](c)
	return &tc
}

type Topic[T tlv.TLVCompatible] struct {
	Chan   *TopicChan[T]
	ConnDb *db.KVStore[*net.Conn]
}

func (topic *Topic[T]) AddConn(conn *net.Conn) {
	k := (*conn).RemoteAddr()
	topic.ConnDb.Set(k.String(), conn)
}

func NewTopic[T tlv.TLVCompatible](name string) *Topic[T] {
	tc := NewTopicChan[T](name)
	topic := &Topic[T]{
		Chan:   tc,
		ConnDb: db.NewKVStore[*net.Conn](nil),
	}
	return topic
}
