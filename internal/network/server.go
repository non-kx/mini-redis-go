package network

import (
	"io"
	"log"
	"net"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

// type RequestHandler func(ctx *RequestContext, arg ...any) error

type IServer interface {
	Start() error
	Stop() error
	HandleConnection(conn net.Conn) error
}

type Server struct {
	Port        string
	Listener    net.Listener
	Connections []net.Conn
	Service     service.IService
	RedisDb     *db.KVStore[[]byte]
	PubsubDb    *db.KVStore[*model.Topic[*tlv.String]]
}

func (s *Server) Start() error {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		log.Println("New connection from:", c.RemoteAddr())
		go s.HandleConnection(c)
	}
}

func (s *Server) Stop() error {
	for _, conn := range s.Connections {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	err := s.Listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) HandleConnection(conn net.Conn) error {
	s.Connections = append(s.Connections, conn)

	for {
		pl := new(payload.RequestPayload)
		_, err := pl.ReadFrom(conn)
		// Handle remove conn from arr
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return err
		}

		reqctx := &payload.RequestContext{
			Conn:     conn,
			Now:      time.Now(),
			Payload:  pl,
			RedisDb:  s.RedisDb,
			PubsubDb: s.PubsubDb,
		}
		err = s.Service.HandleRequest(reqctx)
		if err != nil {
			return err
		}
	}
}

func NewServer(network string, port string, cert string, key string) (*Server, error) {
	l, err := GetListener(network, port, cert, key)
	if err != nil {
		return nil, err
	}

	redisCache := constant.DefaultRedisCachePath
	return &Server{
		Port:        port,
		Listener:    l,
		Connections: make([]net.Conn, 0, 5),
		Service:     service.NewService(),
		RedisDb:     db.NewKVStore[[]byte](&redisCache),
		PubsubDb:    db.NewKVStore[*model.Topic[*tlv.String]](nil),
	}, nil
}
