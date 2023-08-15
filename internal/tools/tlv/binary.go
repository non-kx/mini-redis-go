package tlv

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Binary []byte

func (b *Binary) ToTLV() (TypeLengthValue, error) {
	raw := new(bytes.Buffer)
	b.WriteToIO(raw)

	return TypeLengthValue(raw.Bytes()), nil
}

func (b *Binary) FromTLV(tlv TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	err := b.ReadFromIO(r)
	if err != nil {
		return err
	}

	return nil
}

func (b *Binary) ReadFromIO(r io.Reader) error {
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

	*b = Binary(buf)

	return nil
}

func (b *Binary) WriteToIO(w io.Writer) error {
	typ := BinaryType
	len := len(*b)
	val := *b

	tlv := new(bytes.Buffer)
	binary.Write(tlv, binary.BigEndian, typ)
	binary.Write(tlv, binary.BigEndian, uint32(len))
	binary.Write(tlv, binary.BigEndian, val)

	return nil
}

func (b *Binary) ToString() string {
	return string(*b)
}
