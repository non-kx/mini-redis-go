package payload

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type RequestContext struct {
	context.Context
	Conn     *net.Conn
	RedisDb  *db.KVStore[[]byte]
	PubSubDb *db.KVStore[*Topic[*tlv.String]]
	Now      time.Time
	Payload  *RequestPayload
}

func (ctx *RequestContext) Response(res ResponsePayload) error {
	err := res.WriteToIO(*ctx.Conn)
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

func ReadResponse(r io.Reader) (*ResponsePayload, error) {
	res := new(ResponsePayload)
	err := res.ReadFromIO(r)
	if err != nil {
		return nil, err
	}

	return res, nil
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
