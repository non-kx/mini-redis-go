package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type IServer interface {
	Listen() error
	Close() error
}

type Server struct {
	Port        string
	Listener    net.Listener
	Connections []*net.Conn
}

func (s *Server) Listen() error {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go s.handleNewConn(c)
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

func (s *Server) handleNewConn(conn net.Conn) error {
	s.Connections = append(s.Connections, &conn)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(data)) == "Ping" {
			conn.Write([]byte("Pong\n"))
			continue
		}

		fmt.Print("-> ", string(data))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
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
