package network

import (
	"io"
	"log"
	"net"
	"time"
)

// type RequestHandler func(ctx *RequestContext, arg ...any) error

type IServer interface {
	Start() error
	Stop() error
	HandleConnection(conn *net.Conn) error
	HandleRequest(ctx *RequestContext) error
}

type Server struct {
	Port        string
	Listener    net.Listener
	Connections []*net.Conn
}

func (s *Server) Start() error {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		log.Println("New connection from:", c.RemoteAddr())
		go s.HandleConnection(&c)
	}
}

func (s *Server) Stop() error {
	for _, conn := range s.Connections {
		err := (*conn).Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) HandleConnection(conn *net.Conn) error {
	s.Connections = append(s.Connections, conn)

	for {
		data, err := ReadUntilCRLF(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return err
		}

		reqctx := &RequestContext{
			Now:  time.Now(),
			Data: data,
			Conn: conn,
		}
		err = s.HandleRequest(reqctx)
		if err != nil {
			return err
		}
	}
}

func (s *Server) HandleRequest(ctx *RequestContext) error {
	ctx.Response(ctx.Data)

	return nil
}

func NewServer(network string, port string) (*Server, error) {
	l, err := net.Listen(network, port)
	if err != nil {
		return nil, err
	}

	return &Server{
		Port:        port,
		Listener:    l,
		Connections: make([]*net.Conn, 0, 5),
	}, nil
}
