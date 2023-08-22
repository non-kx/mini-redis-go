package model

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type Topic[T tlv.TLVCompatible] struct {
	Name   string
	ConnDb *db.KVStore[net.Conn]
}

func (topic *Topic[T]) DidInit() bool {
	return topic.ConnDb != nil
}

func (topic *Topic[T]) AddConn(conn net.Conn) {
	k := conn.RemoteAddr()
	topic.ConnDb.Set(k.String(), conn)
}

func (topic *Topic[T]) RemoveConn(conn net.Conn) {
	k := conn.RemoteAddr()
	topic.ConnDb.Delete(k.String())
}

func NewTopic[T tlv.TLVCompatible](name string) *Topic[T] {
	topic := &Topic[T]{
		Name:   name,
		ConnDb: db.NewKVStore[net.Conn](nil),
	}
	return topic
}
