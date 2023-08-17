package tlv

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	EmptyType uint8 = iota
	BinaryType
	StringType

	// Network
	RequestPayloadType
	ResponsePayloadType
	ErrorType

	// Redis
	RedisRequestPayloadType
	GetSuccessType
	SetSuccessType

	// Pubsub
	PubsubRequestPayloadType
	SubSuccessType
	PubSuccessType
	MsgType

	TypeDataLength  uint8  = 1
	LengthDataLegth uint8  = 4
	MaxPayloadSize  uint32 = 10 << 20
)

type TLVCompatible interface {
	io.ReaderFrom
	io.WriterTo
	fmt.Stringer
	ToTLV() (TypeLengthValue, error)
	FromTLV(tlv TypeLengthValue) error
}

type TypeLengthValue []byte

func (tlv *TypeLengthValue) GetType() uint8 {
	if len(*tlv) == 0 {
		return EmptyType
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

func ReadTLVFromIO[T TLVCompatible](r io.Reader) (*T, error) {
	tlv := new(T)
	_, err := (*tlv).ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return tlv, nil
}
