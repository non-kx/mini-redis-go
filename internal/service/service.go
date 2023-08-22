package service

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pingpong"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/pubsub"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type ServiceHandler interface {
	HandleRequest(ctx payload.IRequestContext) error
	HandleDisconnected(conn net.Conn) error
}

type ServiceRequester interface {
	SendPingRequest(conn net.Conn, msg *string) (*tlv.String, error)
	SendGetRequest(conn net.Conn, key string) (tlv.TLVCompatible, error)
	SendSetRequest(conn net.Conn, key string, val tlv.TypeLengthValue) (string, error)
	SendSubRequest(conn net.Conn, topic string) (*pubsub.Subscriber, error)
	SendPubRequest(conn net.Conn, topic string, val tlv.TypeLengthValue) (string, error)
}

type IService interface {
	ServiceHandler
	ServiceRequester
}

type Service struct{}

func (serv *Service) HandleRequest(ctx payload.IRequestContext) error {
	pl := ctx.GetPayload()

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

func (serv *Service) HandleDisconnected(conn net.Conn) error {
	return nil
}

func (serv *Service) SendPingRequest(conn net.Conn, msg *string) (*tlv.String, error) {
	return pingpong.SendPingRequest(conn, msg)
}

func (serv *Service) SendGetRequest(conn net.Conn, key string) (tlv.TLVCompatible, error) {
	return redis.SendGetRequest(conn, key)
}

func (serv *Service) SendSetRequest(conn net.Conn, key string, val tlv.TypeLengthValue) (string, error) {
	return redis.SendSetRequest(conn, key, val)
}

func (serv *Service) SendSubRequest(conn net.Conn, topic string) (*pubsub.Subscriber, error) {
	return pubsub.SendSubRequest(conn, topic)
}

func (serv *Service) SendPubRequest(conn net.Conn, topic string, val tlv.TypeLengthValue) (string, error) {
	return pubsub.SendPubRequest(conn, topic, val)
}

func NewService() *Service {
	return &Service{}
}
