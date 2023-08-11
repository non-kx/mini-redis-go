package pubsub

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

type (
	PubsubCmd     uint8
	PubsubRespTyp uint8
)

const (
	SubscribeCmd PubsubCmd = iota
	PublishCmd
)

const (
	SubscribeSuccess PubsubRespTyp = iota
	PublishSuccess
	Error
)

const (
	CmdDataLength   uint8 = 1
	TopicDataLength uint8 = 32
)

type PubsubRequestPayload struct {
	Cmd     PubsubCmd
	Topic   string
	Message []byte
}

func (payload *PubsubRequestPayload) ToRaw() (*PubsubRawRequestPayload, error) {
	var (
		cmdbyte   []byte
		topicbyte []byte
	)

	cmdbyte = []byte{byte(payload.Cmd)}
	topicbyte = []byte(payload.Topic)
	if len(topicbyte) > int(TopicDataLength) {
		return nil, errors.New("Invalid key length")
	}

	zbyte := make([]byte, int(TopicDataLength)-len(topicbyte))
	topicbyte = append(zbyte, topicbyte...)

	raw := PubsubRawRequestPayload(append(cmdbyte, topicbyte...))

	return &raw, nil
}

type PubsubResponsePayload struct {
	RespType PubsubRespTyp
	RespBody []byte
}

func (payload *PubsubResponsePayload) ToRaw() (*PubsubRawResponsePayload, error) {
	var (
		typbyte  []byte
		bodybyte []byte
	)

	typbyte = []byte(strconv.Itoa(int(payload.RespType)))
	bodybyte = payload.RespBody

	raw := PubsubRawResponsePayload(append(typbyte, bodybyte...))

	return &raw, nil
}

type (
	PubsubRawRequestPayload  []byte
	PubsubRawResponsePayload []byte
)

func (payload *PubsubRawRequestPayload) TransformPayload() (*PubsubRequestPayload, error) {
	var (
		cmd   uint8
		topic = make([]byte, TopicDataLength)
		msg   = make([]byte, len(*payload)-(int(TopicDataLength)+1))
		err   error
	)

	r := bytes.NewReader(*payload)
	err = binary.Read(r, binary.BigEndian, &cmd)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &topic)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &msg)
	if err != nil {
		return nil, err
	}

	return &PubsubRequestPayload{
		Cmd:     PubsubCmd(cmd),
		Topic:   string(topic),
		Message: msg,
	}, nil
}

func (payload *PubsubRawResponsePayload) TransformPayload() (*PubsubResponsePayload, error) {
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

	return &PubsubResponsePayload{
		RespType: PubsubRespTyp(typ),
		RespBody: body,
	}, nil
}
