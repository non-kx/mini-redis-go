package pingpong

import (
	"bytes"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func HandleRequest(ctx *payload.RequestContext) error {
	body := ctx.Payload.Body
	msg := tlv.String("")
	err := msg.ReadFromIO(bytes.NewReader(body))
	if err != nil {
		return err
	}

	var pong tlv.String

	if msg == "PING" {
		pong = "PONG"
	} else {
		pong = msg
	}

	err = pong.WriteToIO(*ctx.Conn)
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	return nil
}

func SendPingRequest(conn *net.Conn, msg *string) (*tlv.String, error) {
	if msg == nil {
		ping := string("PING")
		msg = &ping
	}
	s := tlv.String(*msg)
	msgtlv, err := s.ToTLV()
	if err != nil {
		return nil, err
	}

	req := payload.RequestPayload{
		Cmd:  payload.PingCmd,
		Body: msgtlv,
	}

	err = req.WriteToIO(*conn)
	if err != nil {
		return nil, err
	}

	resp := tlv.String("")
	err = resp.ReadFromIO(*conn)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
