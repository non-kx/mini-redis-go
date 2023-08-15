package pubsub

import (
	"bytes"
	"errors"
	"log"
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func HandleRequest(ctx *payload.RequestContext) error {
	var (
		cmd  uint8
		body tlv.TypeLengthValue
		err  error
	)
	cmd = ctx.Payload.Cmd
	body = ctx.Payload.Body

	pubsubBody := new(payload.PubsubRequestBody)
	err = pubsubBody.ReadFromIO(bytes.NewReader(body))
	if err != nil {
		return err
	}

	switch cmd {
	case payload.SubCmd:
		err = handleSubRequest(ctx, pubsubBody)
		break
	case payload.PubCmd:
		err = handlePubRequest(ctx, pubsubBody)
		break
	default:
		break
	}

	return nil
}

func handleSubRequest(ctx *payload.RequestContext, body *payload.PubsubRequestBody) error {
	var (
		err error
	)
	topic := ctx.PubSubDb.Get(body.Topic)
	if topic == nil {
		topic = payload.NewTopic[*tlv.String](body.Topic)
		ctx.PubSubDb.Set(body.Topic, topic)
	}
	topic.AddConn(ctx.Conn)

	go func() {
		err := handleNewTopic(topic)
		if err != nil {
			log.Panicln(err)
			close(*topic.Chan)
		}
	}()

	err = responseWithSting("OK", ctx)
	if err != nil {
		return err
	}

	return nil
}
func handlePubRequest(ctx *payload.RequestContext, body *payload.PubsubRequestBody) error {
	var (
		err error
	)

	topic := ctx.PubSubDb.Get(body.Topic)
	if topic == nil {
		return errors.New("Error topic does not exist")
	}

	msg := body.Value
	smsg := new(tlv.String)
	err = smsg.FromTLV(msg)
	if err != nil {
		return err
	}

	*topic.Chan <- smsg

	err = responseWithSting("OK", ctx)
	if err != nil {
		return err
	}

	return nil
}

func handleNewTopic[T tlv.TLVCompatible](topic *payload.Topic[T]) error {
	tc := topic.Chan
	conns := topic.ConnDb.Storage

	go func() {
		select {
		case v := <-*tc:
			broadCastMessage[T](v, conns)

		}
	}()
	return nil
}

func broadCastMessage[T tlv.TLVCompatible](msg T, conns map[string]*net.Conn) error {
	for _, conn := range conns {
		c := conn

		raw, err := msg.ToTLV()
		if err != nil {
			return err
		}
		resp := payload.ResponsePayload{
			Typ:  tlv.MsgType,
			Body: raw,
		}
		resp.WriteToIO(*c)
	}
	return nil
}

func responseWithSting(s string, ctx *payload.RequestContext) error {
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
	err = resp.WriteToIO(*ctx.Conn)
	if err != nil {
		err = ctx.Error(uint16(tlv.DataTransformError), tlv.ErrMsg[tlv.DataTransformError])
		log.Println(err)
		return err
	}

	return nil
}
