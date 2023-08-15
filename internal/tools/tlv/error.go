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

func (e *Error) ReadFromIO(r io.Reader) error {
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

	if typ != ErrorType {
		return errors.New("Invalid Error type")
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

	*e = Error{
		Code: ErrCode(binary.BigEndian.Uint16(buf[0:2])),
		Msg:  string(buf[2:]),
	}
	return nil
}

func (e *Error) WriteToIO(w io.Writer) error {
	var (
		typ        = ErrorType
		len uint32 = uint32(2 + len(e.Msg))
	)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, e.Code)
	binary.Write(buf, binary.BigEndian, e.Msg)

	val := buf.Bytes()

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, len)
	binary.Write(w, binary.BigEndian, val)

	return nil
}

func NewError(code uint16, msg string) Error {
	return Error{
		Code: ErrCode(code),
		Msg:  msg,
	}
}
