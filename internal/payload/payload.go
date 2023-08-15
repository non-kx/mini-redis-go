package payload

import (
	"bytes"
	"encoding/binary"
	"io"

	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

// type (
// 	Command uint8
// 	RespTyp uint8
// )

const (
	PingResp uint8 = iota
	GetResp
	SetResp
	SubResp
	PubResp

	ErrResp = 99
)

const (
	PingCmd uint8 = iota
	GetCmd
	SetCmd
	SubCmd
	PubCmd
)

type RawRequestPayload []byte
type RequestPayload struct {
	Cmd  uint8
	Body tlv.TypeLengthValue
}

func (req *RequestPayload) ReadFromIO(r io.Reader) error {
	var (
		typ uint8
		len uint32
		buf []byte
		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return err
	}

	buf = make([]byte, len)
	_, err = r.Read(buf)
	if err != nil {
		return err
	}

	*req = RequestPayload{
		Cmd:  buf[0],
		Body: buf[1:],
	}
	return nil
}

func (req *RequestPayload) WriteToIO(w io.Writer) error {
	var (
		typ = tlv.RequestPayloadType
		val []byte
	)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, req.Cmd)
	binary.Write(buf, binary.BigEndian, req.Body)

	val = buf.Bytes()

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, uint32(len(val)))
	binary.Write(w, binary.BigEndian, val)

	return nil
}

type RawResponsePayload []byte
type ResponsePayload struct {
	Typ  uint8
	Body tlv.TypeLengthValue
}

func (res *ResponsePayload) ReadFromIO(r io.Reader) error {
	var (
		typ uint8
		len uint32
		buf []byte
		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return err
	}

	buf = make([]byte, len)
	_, err = r.Read(buf)
	if err != nil {
		return err
	}

	*res = ResponsePayload{
		Typ:  buf[0],
		Body: buf[1:],
	}
	return nil
}

func (res *ResponsePayload) WriteToIO(w io.Writer) error {
	var (
		typ = tlv.RequestPayloadType
		val []byte
	)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, res.Typ)
	binary.Write(buf, binary.BigEndian, res.Body)

	val = buf.Bytes()

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, uint32(len(val)))
	binary.Write(w, binary.BigEndian, val)

	return nil
}
