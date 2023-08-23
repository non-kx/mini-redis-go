package redis

import (
	"bytes"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func HandleRequest(ctx payload.IRequestContext) error {
	var (
		cmd  uint8
		body tlv.TypeLengthValue
		err  error
	)
	pl := ctx.GetPayload()
	cmd = pl.Cmd
	body = pl.Body

	redisBody := new(payload.RedisRequestBody)
	_, err = redisBody.ReadFrom(bytes.NewReader(body))
	if err != nil {
		return err
	}

	switch cmd {
	case payload.GetCmd:
		err = handleGetRequest(ctx, redisBody)
		break
	case payload.SetCmd:
		err = handleSetRequest(ctx, redisBody)
		break
	default:
		break
	}

	return nil
}

func handleGetRequest(ctx payload.IRequestContext, body *payload.RedisRequestBody) error {
	var (
		raw tlv.TypeLengthValue
		typ uint8
		err error
	)
	data := ctx.GetRedis(body.Key)
	raw = tlv.TypeLengthValue(data)
	typ = raw.GetType()

	resp := payload.ResponsePayload{
		Typ:  typ,
		Body: raw,
	}
	_, err = resp.WriteTo(ctx.GetConn())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func handleSetRequest(ctx payload.IRequestContext, body *payload.RedisRequestBody) error {
	var (
		resp payload.ResponsePayload
		err  error
	)
	ctx.SetRedis(body.Key, body.Value)

	s := tlv.String("OK")
	raw, err := s.ToTLV()
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	resp = payload.ResponsePayload{
		Typ:  tlv.StringType,
		Body: raw,
	}
	_, err = resp.WriteTo(ctx.GetConn())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
