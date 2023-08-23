package pubsub

import (
	"bytes"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/db/model"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/service/internal/helper"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"bitbucket.org/non-pn/mini-redis-go/pkg/async"
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

	pubsubBody := new(payload.PubsubRequestBody)
	_, err = pubsubBody.ReadFrom(bytes.NewReader(body))
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

func handleSubRequest(ctx payload.IRequestContext, body *payload.PubsubRequestBody) error {
	var (
		err error
	)
	topic := ctx.GetPubsub(body.Topic)
	if topic == nil || !(*topic).DidInit() {
		topic = model.NewTopic[*tlv.String](body.Topic)
		ctx.SetPubsub(body.Topic, topic)
	}
	(*topic).AddConn(ctx.GetConn())

	err = helper.ResponseWithString("OK", ctx)
	if err != nil {
		return err
	}

	return nil
}
func handlePubRequest(ctx payload.IRequestContext, body *payload.PubsubRequestBody) error {
	var (
		err error
	)

	topic := ctx.GetPubsub(body.Topic)
	if topic == nil || !(*topic).DidInit() {
		topic = model.NewTopic[*tlv.String](body.Topic)
		ctx.SetPubsub(body.Topic, topic)
	}

	msg := body.Value
	smsg := new(tlv.String)
	err = smsg.FromTLV(msg)
	if err != nil {
		return err
	}

	rc, ec := async.Async(func() (any, error) {
		err := broadCastMessage[*tlv.String](topic, smsg)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	_, err = async.Await(rc, ec)
	if err != nil {
		return err
	}

	err = helper.ResponseWithString("OK", ctx)
	if err != nil {
		return err
	}

	return nil
}

func broadCastMessage[T tlv.TLVCompatible](topic *model.Topic[T], msg T) error {
	conns := topic.ConnDb.Storage
	for _, conn := range conns {
		c := conn

		log.Println("Relaying message to:", c.RemoteAddr().String())

		raw, err := msg.ToTLV()
		if err != nil {
			return err
		}
		resp := payload.ResponsePayload{
			Typ:  tlv.MsgType,
			Body: raw,
		}
		resp.WriteTo(c)
	}
	return nil
}
