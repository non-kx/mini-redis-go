package pubsub

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

type Subscriber struct {
	Conn     *net.Conn
	Messages []*tlv.String
}

func (sub *Subscriber) NextMessage() (*tlv.String, error) {
	// if len(sub.Messages) == 0 {
	// 	return nil
	// }

	// msg := sub.Messages[0]
	// sub.Messages = sub.Messages[1:]

	// return msg
	pl := new(payload.ResponsePayload)
	err := pl.ReadFromIO(*sub.Conn)
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

func (sub *Subscriber) Subscribe() error {
	for {
		pl := new(payload.ResponsePayload)
		err := pl.ReadFromIO(*sub.Conn)
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

		fmt.Println("new msg:", s)

		sub.Messages = append(sub.Messages, s)
	}
}

func NewSubscriber(conn *net.Conn) *Subscriber {
	sub := &Subscriber{
		Conn:     conn,
		Messages: make([]*tlv.String, 0),
	}

	return sub
}
