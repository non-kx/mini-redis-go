package service

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pingpong"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pubsub"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
)

type Handler interface {
	HandleRequest(ctx *payload.RequestContext) error
}

func HandleRequest(ctx *payload.RequestContext) error {
	pl := ctx.Payload

	switch pl.Cmd {
	case payload.PingCmd:
		err := pingpong.HandleRequest(ctx)
		if err != nil {
			return err
		}
		break
	case payload.GetCmd:
		err := redis.HandleRequest(ctx)
		if err != nil {
			return err
		}
		break
	case payload.SetCmd:
		err := redis.HandleRequest(ctx)
		if err != nil {
			return err
		}
		break
	case payload.SubCmd:
		err := pubsub.HandleRequest(ctx)
		if err != nil {
			return err
		}
		break
	case payload.PubCmd:
		err := pubsub.HandleRequest(ctx)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func HandleDisconnected(conn *net.Conn) error {
	return nil
}
