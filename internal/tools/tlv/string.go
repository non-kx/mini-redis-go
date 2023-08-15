package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type String string

func (s *String) ReadFromIO(r io.Reader) error {
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

	if typ != StringType {
		return errors.New("Invalid String")
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

	*s = String(buf)

	return nil
}

func (s *String) WriteToIO(w io.Writer) error {
	var (
		typ        = StringType
		len uint32 = uint32(len(*s))
	)

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, len)
	binary.Write(w, binary.BigEndian, []byte(*s))

	return nil
}

func (s *String) ToTLV() (TypeLengthValue, error) {
	typ := StringType
	len := len(*s)
	val := []byte(*s)

	tlv := new(bytes.Buffer)
	binary.Write(tlv, binary.BigEndian, typ)
	binary.Write(tlv, binary.BigEndian, uint32(len))
	binary.Write(tlv, binary.BigEndian, val)

	return TypeLengthValue(tlv.Bytes()), nil
}

func (s *String) FromTLV(tlv TypeLengthValue) error {
	var (
		typ uint8
		len uint32
		buf []byte
		err error
	)

	r := bytes.NewReader(tlv)
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

	*s = String(buf)

	return nil
}

func (s *String) ToString() string {
	if s == nil {
		return ""
	}
	return string(*s)
}
