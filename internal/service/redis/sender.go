package redis

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func SendGetRequest(conn *net.Conn, key string) (tlv.TLVCompatible, error) {
	body := payload.RedisRequestBody{
		Key:   key,
		Len:   0,
		Value: []byte{},
	}
	rawbod, err := body.ToTLV()
	if err != nil {
		return nil, err
	}

	req := payload.RequestPayload{
		Cmd:  payload.GetCmd,
		Body: rawbod,
	}

	_, err = req.WriteTo(*conn)
	if err != nil {
		return nil, err
	}

	resp := new(payload.ResponsePayload)
	_, err = resp.ReadFrom(*conn)
	if err != nil {
		return nil, err
	}

	var val tlv.TLVCompatible
	switch resp.Typ {
	case tlv.StringType:
		val = new(tlv.String)
		err = val.FromTLV(resp.Body)
		if err != nil {
			return nil, err
		}
		break
	case tlv.BinaryType:
		val = new(tlv.Binary)
		err = val.FromTLV(resp.Body)
		if err != nil {
			return nil, err
		}
		break
	default:
		break
	}

	return val, nil
}

func SendSetRequest(conn *net.Conn, key string, val tlv.TypeLengthValue) (string, error) {
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

	_, err = req.WriteTo(*conn)
	if err != nil {
		return "", err
	}

	resp := new(payload.ResponsePayload)
	_, err = resp.ReadFrom(*conn)
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
