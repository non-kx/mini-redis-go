package payload

import (
	"context"
	"net"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type StringTopic *model.Topic[*tlv.String]

type IRequestContext interface {
	Response(res ResponsePayload) error
	Error(code uint16, msg string) error
	GetConn() net.Conn
	GetPayload() *RequestPayload
	GetRedis(k string) []byte
	SetRedis(k string, v []byte)
	GetPubsub(k string) StringTopic
	SetPubsub(k string, v StringTopic)
}

type RequestContext struct {
	context.Context
	Conn     net.Conn
	RedisDb  *db.KVStore[[]byte]
	PubsubDb *db.KVStore[*model.Topic[*tlv.String]]
	Now      time.Time
	Payload  *RequestPayload
}

func (ctx *RequestContext) Response(res ResponsePayload) error {
	_, err := res.WriteTo(ctx.Conn)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *RequestContext) Error(code uint16, msg string) error {
	tlvErr := tlv.NewError(code, msg)
	_, err := tlvErr.WriteTo(ctx.Conn)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *RequestContext) GetConn() net.Conn {
	return ctx.Conn
}

func (ctx *RequestContext) GetPayload() *RequestPayload {
	return ctx.Payload
}

func (ctx *RequestContext) GetRedis(k string) []byte {
	return ctx.RedisDb.Get(k)
}
func (ctx *RequestContext) SetRedis(k string, v []byte) {
	ctx.RedisDb.Set(k, v)
}
func (ctx *RequestContext) GetPubsub(k string) StringTopic {
	return ctx.PubsubDb.Get(k)
}
func (ctx *RequestContext) SetPubsub(k string, v StringTopic) {
	ctx.PubsubDb.Set(k, v)
}
