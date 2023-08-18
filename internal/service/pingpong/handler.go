package pingpong

import (
	"bytes"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func HandleRequest(ctx *payload.RequestContext) error {
	body := ctx.Payload.Body
	msg := tlv.String("")
	_, err := msg.ReadFrom(bytes.NewReader(body))
	if err != nil {
		return err
	}

	var pong tlv.String

	if msg == "PING" {
		pong = "PONG"
	} else {
		pong = msg
	}

	_, err = pong.WriteTo(*ctx.Conn)
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	return nil
}
