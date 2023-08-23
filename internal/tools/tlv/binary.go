package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Binary []byte

func (b *Binary) ReadFrom(r io.Reader) (int64, error) {
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

	if typ != BinaryType {
		return n, errors.New("Invalid Binary")
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
	*b = Binary(buf)

	return n, nil
}

func (b *Binary) WriteTo(w io.Writer) (int64, error) {
	var (
		typ        = BinaryType
		len uint32 = uint32(len(*b))
		n   int64

		err error
	)

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, uint32(len))
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(*b))
	if err != nil {
		return n, err
	}

	n += int64(o)

	return n, err
}

func (b *Binary) FromTLV(tlv TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := b.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}

func (b *Binary) ToTLV() (TypeLengthValue, error) {
	tlv := new(bytes.Buffer)
	_, err := b.WriteTo(tlv)
	if err != nil {
		return nil, err
	}

	return TypeLengthValue(tlv.Bytes()), nil
}

func (b *Binary) String() string {
	return string(*b)
}
