package payload

import (
	"context"
	"net"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type IRequestContext interface {
	Response(res ResponsePayload) error
	Error(code uint16, msg string) error
}

type RequestContext struct {
	context.Context
	Conn     *net.Conn
	RedisDb  *db.KVStore[[]byte]
	PubSubDb *db.KVStore[*model.Topic[*tlv.String]]
	Now      time.Time
	Payload  *RequestPayload
}

func (ctx *RequestContext) Response(res ResponsePayload) error {
	_, err := res.WriteTo(*ctx.Conn)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *RequestContext) Error(code uint16, msg string) error {
	tlvErr := tlv.NewError(code, msg)
	err := tlvErr.WriteToIO(*ctx.Conn)
	if err != nil {
		return err
	}

	return nil
}
