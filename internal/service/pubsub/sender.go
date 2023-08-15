package pubsub

import (
	"fmt"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func SendSubRequest(conn *net.Conn, topic string) (*Subscriber, error) {
	body := payload.PubsubRequestBody{
		Topic: topic,
		Len:   0,
		Value: []byte{},
	}
	rawbod, err := body.ToTLV()
	if err != nil {
		return nil, err
	}

	req := payload.RequestPayload{
		Cmd:  payload.SubCmd,
		Body: rawbod,
	}

	err = req.WriteToIO(*conn)
	if err != nil {
		return nil, err
	}

	resp := new(payload.ResponsePayload)
	err = resp.ReadFromIO(*conn)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.Body)

	sub := NewSubscriber(conn)

	return sub, nil
}

func SendPubRequest(conn *net.Conn, key string, val tlv.TypeLengthValue) (string, error) {
	body := payload.RedisRequestBody{
		Key:   key,
		Len:   uint32(len(val)),
		Value: val,
	}
	rawbod, err := body.ToTLV()
	if err != nil {
		return "", err
	}

	req := payload.RequestPayload{
		Cmd:  payload.SetCmd,
		Body: rawbod,
	}

	err = req.WriteToIO(*conn)
	if err != nil {
		return "", err
	}

	resp := new(payload.ResponsePayload)
	err = resp.ReadFromIO(*conn)
	if err != nil {
		return "", err
	}

	msg := new(tlv.String)
	err = msg.FromTLV(resp.Body)
	if err != nil {
		return "", err
	}

	return string(*msg), nil
}
