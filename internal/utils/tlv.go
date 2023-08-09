package utils

import (
	"bytes"
	"encoding/binary"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20
)

type TypeLengthValue []byte

func (tlv *TypeLengthValue) GetType() uint8 {
	return (*tlv)[0]
}

func (tlv *TypeLengthValue) GetLength() uint32 {
	return binary.BigEndian.Uint32((*tlv)[1:5])
}

func (tlv *TypeLengthValue) GetValue() []byte {
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
	val := *s

	tvl := new(bytes.Buffer)
	binary.Write(tvl, binary.BigEndian, typ)
	binary.Write(tvl, binary.BigEndian, uint32(len))
	binary.Write(tvl, binary.BigEndian, val)

	return TypeLengthValue(tvl.Bytes()), nil
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

	tvl := new(bytes.Buffer)
	binary.Write(tvl, binary.BigEndian, typ)
	binary.Write(tvl, binary.BigEndian, uint32(len))
	binary.Write(tvl, binary.BigEndian, val)

	return TypeLengthValue(tvl.Bytes()), nil
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

	*b = Binary(buf)

	return nil
}
