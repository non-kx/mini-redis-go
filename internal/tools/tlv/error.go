package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type ErrCode uint16

const (
	DataTransformError ErrCode = iota
)

var ErrMsg = map[ErrCode]string{
	DataTransformError: "Data transform error",
}

type Error struct {
	Code ErrCode
	Msg  string
}

func (e *Error) ReadFrom(r io.Reader) (int64, error) {
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

	if typ != ErrorType {
		return n, errors.New("Invalid Error type")
	}

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

	*e = Error{
		Code: ErrCode(binary.BigEndian.Uint16(buf[0:2])),
		Msg:  string(buf[2:]),
	}
	return n, nil
}

func (e *Error) WriteTo(w io.Writer) (int64, error) {
	var (
		typ        = ErrorType
		len uint32 = uint32(2 + len([]byte(e.Msg)))
		n   int64

		err error
	)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, e.Code)
	if err != nil {
		return 0, err
	}
	err = binary.Write(buf, binary.BigEndian, []byte(e.Msg))
	if err != nil {
		return 0, err
	}

	val := buf.Bytes()

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	if typ != ErrorType {
		return n, errors.New("Invalid Error type")
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, len)
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

func (e *Error) FromTLV(tlv TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := e.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}

func (e *Error) ToTLV() (TypeLengthValue, error) {
	tlv := new(bytes.Buffer)
	_, err := e.WriteTo(tlv)
	if err != nil {
		return nil, err
	}

	return TypeLengthValue(tlv.Bytes()), nil
}

func NewError(code uint16, msg string) Error {
	return Error{
		Code: ErrCode(code),
		Msg:  msg,
	}
}
