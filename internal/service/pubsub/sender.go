package pubsub

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func SendSubRequest(conn *net.Conn, topic string) (*Subscriber, error) {
	body := payload.PubsubRequestBody{
		Topic: topic,
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

	_, err = req.WriteTo(*conn)
	if err != nil {
		return nil, err
	}

	_, err = payload.ReadResponse(*conn)
	if err != nil {
		return nil, err
	}

	sub := NewSubscriber(conn)

	return sub, nil
}

func SendPubRequest(conn *net.Conn, topic string, val tlv.TypeLengthValue) (string, error) {
	body := payload.PubsubRequestBody{
		Topic: topic,
		Value: val,
	}
	rawbod, err := body.ToTLV()
	if err != nil {
		return "", err
	}

	req := payload.RequestPayload{
		Cmd:  payload.PubCmd,
		Body: rawbod,
	}

	_, err = req.WriteTo(*conn)
	if err != nil {
		return "", err
	}

	resp, err := payload.ReadResponse(*conn)
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
