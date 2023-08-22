package pubsub

import (
	"errors"
	"io"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type Subscriber struct {
	Conn         net.Conn
	IsSubscribed bool
}

func (sub *Subscriber) NextMessage() (*tlv.String, error) {
	if !sub.IsSubscribed {
		pl := new(payload.ResponsePayload)
		_, err := pl.ReadFrom(sub.Conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return nil, err
		}

		if pl.Typ != tlv.MsgType {
			return nil, errors.New("Error not message type")
		}

		s := new(tlv.String)
		err = s.FromTLV(pl.Body)
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	return nil, nil
}

func (sub *Subscriber) Subscribe(handler func(string)) error {
	sub.IsSubscribed = true
	defer func() {
		sub.IsSubscribed = false
	}()
	for {
		pl := new(payload.ResponsePayload)
		_, err := pl.ReadFrom(sub.Conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return err
		}

		if pl.Typ != tlv.MsgType {
			return errors.New("Error not message type")
		}

		s := new(tlv.String)
		err = s.FromTLV(pl.Body)
		if err != nil {
			return err
		}

		handler(s.String())
	}
}

func NewSubscriber(conn net.Conn) *Subscriber {
	sub := &Subscriber{
		Conn:         conn,
		IsSubscribed: false,
	}

	return sub
}
