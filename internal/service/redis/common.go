package redis

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	"bitbucket.org/non-pn/mini-redis-go/internal/utils"
)

type (
	RedisCmd      uint8
	RedisRespType uint8
)

const (
	PROTOCOL = "tcp"
)

const (
	GetCmd RedisCmd = iota
	SetCmd
)

const (
	GetSuccess RedisRespType = iota
	SetSuccess
	Error
)

const (
	CmdDataLength uint8 = 1
	KeyDataLength uint8 = 16
)

type RedisRequestPayload struct {
	Cmd   RedisCmd
	Key   string
	Value utils.TypeLengthValue
}

func (payload *RedisRequestPayload) ToRaw() (*RedisRawRequestPayload, error) {
	var (
		cmdbyte []byte
		kbyte   []byte
		vbyte   []byte
	)

	cmdbyte = []byte{byte(payload.Cmd)}
	kbyte = []byte(payload.Key)
	if len(kbyte) > int(KeyDataLength) {
		return nil, errors.New("Invalid key length")
	}

	zbyte := make([]byte, int(KeyDataLength)-len(kbyte))
	kbyte = append(zbyte, kbyte...)
	vbyte = []byte(payload.Value)

	raw := RedisRawRequestPayload(append(append(cmdbyte, kbyte...), vbyte...))

	return &raw, nil
}

type RedisResponsePayload struct {
	RespType RedisRespType
	RespBody []byte
}

func (payload *RedisResponsePayload) ToRaw() (*RedisRawResponsePayload, error) {
	var (
		typbyte  []byte
		bodybyte []byte
	)

	typbyte = []byte(strconv.Itoa(int(payload.RespType)))
	bodybyte = payload.RespBody

	raw := RedisRawResponsePayload(append(typbyte, bodybyte...))

	return &raw, nil
}

type (
	RedisRawRequestPayload  []byte
	RedisRawResponsePayload []byte
)

func (payload *RedisRawRequestPayload) TransformPayload() (*RedisRequestPayload, error) {
	var (
		cmd uint8
		key = make([]byte, KeyDataLength)
		val = make([]byte, len(*payload)-(1+int(KeyDataLength)))
		err error
	)

	r := bytes.NewReader(*payload)
	err = binary.Read(r, binary.BigEndian, &cmd)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &key)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &val)
	if err != nil {
		return nil, err
	}

	tlv := utils.TypeLengthValue(val)

	return &RedisRequestPayload{
		Cmd:   RedisCmd(cmd),
		Key:   string(key),
		Value: tlv,
	}, nil
}

func (payload *RedisRawResponsePayload) TransformPayload() (*RedisResponsePayload, error) {
	var (
		typ  uint8
		body = make([]byte, len(*payload)-1)
		err  error
	)

	r := bytes.NewReader(*payload)
	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body)
	if err != nil {
		return nil, err
	}

	return &RedisResponsePayload{
		RespType: RedisRespType(typ),
		RespBody: body,
	}, nil
}
