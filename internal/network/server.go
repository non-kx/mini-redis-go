package network

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type RequestHandler func(ctx *RequestContext) error

type IServer interface {
	Listen() error
	Close() error
}

type Server struct {
	Port        string
	Listener    net.Listener
	Connections []*net.Conn
	Handler     RequestHandler
}

func (s *Server) Listen() error {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		log.Println("New connection from:", c.RemoteAddr())
		go s.HandleConnection(&c)
	}
}

func (s *Server) Close() error {
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
		data, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return err
		}

		reqctx := &RequestContext{
			Now:  time.Now(),
			Data: []byte(strings.TrimSuffix(data, "\n")),
			conn: conn,
		}
		err = s.Handler(reqctx)
		if err != nil {
			return err
		}
	}
}

func NewServer(network string, port string, handler RequestHandler) (*Server, error) {
	l, err := net.Listen(network, port)
	if err != nil {
		return nil, err
	}

	return &Server{
		Port:        port,
		Listener:    l,
		Connections: make([]*net.Conn, 0, 5),
		Handler:     handler,
	}, nil
}
