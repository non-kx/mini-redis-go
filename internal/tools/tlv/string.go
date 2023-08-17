package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type String string

func (s *String) ReadFrom(r io.Reader) (int64, error) {
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

	if typ != StringType {
		return n, errors.New("Invalid String")
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

	*s = String(buf)

	return n, nil
}

func (s *String) WriteTo(w io.Writer) (int64, error) {
	var (
		typ        = StringType
		len uint32 = uint32(len(*s))
		n   int64

		err error
	)

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, len)
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(*s))

	n += int64(o)

	return n, err
}

func (s *String) FromTLV(tlv TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := s.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}

func (s *String) ToTLV() (TypeLengthValue, error) {
	tlv := new(bytes.Buffer)
	_, err := s.WriteTo(tlv)
	if err != nil {
		return nil, err
	}

	return TypeLengthValue(tlv.Bytes()), nil
}

func (s *String) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}
