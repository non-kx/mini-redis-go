package network

import (
	"context"
	"net"
	"time"
)

type RequestContext struct {
	context.Context
	conn *net.Conn
	Now  time.Time
	Data []byte
}

func (ctx *RequestContext) Response(data []byte) error {
	_, err := (*ctx.conn).Write(append(data, []byte("\n")...))
	if err != nil {
		return err
	}

	return nil
}
