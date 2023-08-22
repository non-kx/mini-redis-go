package helper

import (
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func ResponseWithString(s string, ctx payload.IRequestContext) error {
	ss := tlv.String(s)
	raw, err := ss.ToTLV()
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	typ := raw.GetType()
	resp := payload.ResponsePayload{
		Typ:  typ,
		Body: raw,
	}
	_, err = resp.WriteTo(ctx.GetConn())
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	return nil
}
