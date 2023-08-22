package payload

import (
	"bytes"
	"encoding/binary"
	"io"

	tlvpac "bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
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
	Body tlvpac.TypeLengthValue
}

func (req *RequestPayload) ReadFrom(r io.Reader) (int64, error) {
	var (
		typ uint8
		len uint32
		buf []byte
		n   int64

		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return n, err
	}

	n += 4

	buf = make([]byte, len)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	n += int64(o)

	*req = RequestPayload{
		Cmd:  buf[0],
		Body: buf[1:],
	}
	return n, nil
}

func (req *RequestPayload) WriteTo(w io.Writer) (int64, error) {
	var (
		typ = tlvpac.RequestPayloadType
		val []byte
		n   int64

		err error
	)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, req.Cmd)
	if err != nil {
		return 0, err
	}
	err = binary.Write(buf, binary.BigEndian, req.Body)
	if err != nil {
		return 0, err
	}

	val = buf.Bytes()

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	n += 1

	binary.Write(w, binary.BigEndian, uint32(len(val)))
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(val)
	if err != nil {
		return n, err
	}

	n += int64(o)

	return n, nil
}

func (req *RequestPayload) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := req.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}

func (req *RequestPayload) ToTLV() (tlvpac.TypeLengthValue, error) {
	tlv := new(bytes.Buffer)
	_, err := req.WriteTo(tlv)
	if err != nil {
		return nil, err
	}

	return tlvpac.TypeLengthValue(tlv.Bytes()), nil
}

type RawResponsePayload []byte
type ResponsePayload struct {
	Typ  uint8
	Body tlvpac.TypeLengthValue
}

func (res *ResponsePayload) ReadFrom(r io.Reader) (int64, error) {
	var (
		typ uint8
		len uint32
		buf []byte
		n   int64

		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return n, err
	}

	n += 4

	buf = make([]byte, len)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	n += int64(o)

	*res = ResponsePayload{
		Typ:  buf[0],
		Body: buf[1:],
	}
	return n, nil
}

func (res *ResponsePayload) WriteTo(w io.Writer) (int64, error) {
	var (
		typ = tlvpac.RequestPayloadType
		val []byte
		n   int64

		err error
	)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, res.Typ)
	if err != nil {
		return 0, err
	}
	binary.Write(buf, binary.BigEndian, res.Body)
	if err != nil {
		return 0, err
	}

	val = buf.Bytes()

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, uint32(len(val)))
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(val)
	if err != nil {
		return n, err
	}

	n += int64(o)

	return n, nil
}

func (res *ResponsePayload) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := res.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}

func (res *ResponsePayload) ToTLV() (tlvpac.TypeLengthValue, error) {
	tlv := new(bytes.Buffer)
	_, err := res.WriteTo(tlv)
	if err != nil {
		return nil, err
	}

	return tlvpac.TypeLengthValue(tlv.Bytes()), nil
}

func ReadResponse(r io.Reader) (*ResponsePayload, error) {
	res := new(ResponsePayload)
	_, err := res.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
