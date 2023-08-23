package test

import (
	"testing"

	mocknet "bitbucket.org/non-pn/mini-redis-go/internal/mock/net"
	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func ExpectWriteRequestToConn(t *testing.T, conn *mocknet.MockConn, cmd uint8, body tlv.TypeLengthValue) {
	req := payload.RequestPayload{
		Cmd:  cmd,
		Body: body,
	}
	raw, err := req.ToTLV()
	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Write(typ).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Write(blen).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Write(val).SetArg(0, val).Return(len(val), nil)
}

func ExpectReadResponseFromConn(t *testing.T, conn *mocknet.MockConn, restyp uint8, body tlv.TypeLengthValue) {
	resp := payload.ResponsePayload{
		Typ:  restyp,
		Body: body,
	}
	raw, err := resp.ToTLV()
	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Read(gomock.Any()).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, val).Return(len(val), nil)
}

func ExpectReadStringResponseFromConn(t *testing.T, conn *mocknet.MockConn, str string) {
	s := tlv.String(str)
	raw, err := s.ToTLV()

	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Read(gomock.Any()).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Read(gomock.Any()).SetArg(0, val).Return(len(val), nil)
}

func ExpectWriteResponseToConn(t *testing.T, conn *mocknet.MockConn, restyp uint8, body tlv.TypeLengthValue) {
	resp := payload.ResponsePayload{
		Typ:  restyp,
		Body: body,
	}
	raw, err := resp.ToTLV()
	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Write(typ).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Write(blen).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Write(val).SetArg(0, val).Return(len(val), nil)
}

func ExpectWriteStringResponseToConn(t *testing.T, conn *mocknet.MockConn, str string) {
	s := tlv.String(str)
	raw, err := s.ToTLV()

	assert.Nil(t, err)

	typ := raw[0:1]
	blen := raw[1:5]
	val := raw[5:]

	conn.EXPECT().Write(typ).SetArg(0, typ).Return(len(typ), nil)
	conn.EXPECT().Write(blen).SetArg(0, blen).Return(len(blen), nil)
	conn.EXPECT().Write(val).SetArg(0, val).Return(len(val), nil)
}
