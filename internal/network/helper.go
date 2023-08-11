package network

import (
	"bufio"
	"context"
	"net"
	"time"
)

type RequestContext struct {
	context.Context
	Conn *net.Conn
	Now  time.Time
	Data []byte
}

func (ctx *RequestContext) Response(data []byte) error {
	err := WriteWithCRLF(ctx.Conn, data)
	if err != nil {
		return err
	}

	return nil
}

// dropCR drops a terminal \r\n from the data.
func dropCRLF(data []byte) []byte {
	if len(data) > 1 {
		if string(data[len(data)-2:]) == "\r\n" {
			return data[0 : len(data)-2]
		}
	}
	return data
}

func ReadUntilCRLF(conn *net.Conn) ([]byte, error) {
	data, err := bufio.NewReader(*conn).ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return dropCRLF(data), nil
}

func WriteWithCRLF(conn *net.Conn, data []byte) error {
	_, err := (*conn).Write(append(data, []byte("\r\n")...))
	if err != nil {
		return err
	}

	return nil
}
