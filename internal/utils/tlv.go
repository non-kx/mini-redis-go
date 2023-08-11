package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	NilType uint8 = iota
	BinaryType
	StringType

	MaxPayloadSize uint32 = 10 << 20
)

type TypeLengthValue []byte

func (tlv *TypeLengthValue) GetType() uint8 {
	if len(*tlv) == 0 {
		return NilType
	}
	return (*tlv)[0]
}

func (tlv *TypeLengthValue) GetLength() uint32 {
	if len(*tlv) < 5 {
		return 0
	}
	return binary.BigEndian.Uint32((*tlv)[1:5])
}

func (tlv *TypeLengthValue) GetValue() []byte {
	if len(*tlv) < 6 {
		return []byte{}
	}
	return (*tlv)[5:]
}

type TLVCompatible interface {
	ToTLV() (TypeLengthValue, error)
	FromTLV(tlv TypeLengthValue) error
}

type String string

func (s *String) ToTLV() (TypeLengthValue, error) {
	typ := StringType
	len := len(*s)
	val := []byte(*s)

	tlv := new(bytes.Buffer)
	binary.Write(tlv, binary.BigEndian, typ)
	binary.Write(tlv, binary.BigEndian, uint32(len))
	binary.Write(tlv, binary.BigEndian, val)

	fmt.Println(tlv.Bytes())

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

type Binary []byte

func (b *Binary) ToTLV() (TypeLengthValue, error) {
	typ := BinaryType
	len := len(*b)
	val := *b

	tlv := new(bytes.Buffer)
	binary.Write(tlv, binary.BigEndian, typ)
	binary.Write(tlv, binary.BigEndian, uint32(len))
	binary.Write(tlv, binary.BigEndian, val)

	return TypeLengthValue(tlv.Bytes()), nil
}

func (b *Binary) FromTLV(tlv TypeLengthValue) error {
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

	fmt.Println(typ, len, buf)

	*b = Binary(buf)

	return nil
}
