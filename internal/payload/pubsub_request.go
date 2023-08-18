package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	tlvpac "bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

const (
	TopicDataLength uint8 = 16
)

type PubsubRequestBody struct {
	Topic string
	Value tlvpac.TypeLengthValue
}

func (b *PubsubRequestBody) ReadFrom(r io.Reader) (int64, error) {
	var (
		typ   uint8
		len   uint32
		topic = make([]byte, TopicDataLength)
		n     int64

		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	n += 1

	if typ != tlvpac.PubsubRequestPayloadType {
		return n, errors.New("Error: Invalid pubsub request")
	}

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return n, err
	}

	n += 4

	err = binary.Read(r, binary.BigEndian, &topic)
	if err != nil {
		return n, err
	}

	n += int64(TopicDataLength)

	vlen := len - uint32(TopicDataLength)
	buf := make([]byte, vlen)

	if vlen > 0 {
		o, err := r.Read(buf)
		if err != nil {
			return n, err
		}

		n += int64(o)
	}

	*b = PubsubRequestBody{
		Topic: string(topic),
		Value: tlvpac.TypeLengthValue(buf),
	}

	return n, nil
}

func (b *PubsubRequestBody) WriteTo(w io.Writer) (int64, error) {
	var (
		n int64

		err error
	)

	if len(b.Topic) > int(TopicDataLength) {
		return 0, errors.New("Error: Topic length exceed limit")
	}

	typ := tlvpac.PubsubRequestPayloadType
	topic := append(make([]byte, int(TopicDataLength)-len(b.Topic)), []byte(b.Topic)...)
	val := b.Value
	blen := uint32(len(val)) + uint32(TopicDataLength)

	err = binary.Write(w, binary.BigEndian, typ)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, blen)
	if err != nil {
		return n, err
	}

	n += 4

	err = binary.Write(w, binary.BigEndian, topic)
	if err != nil {
		return n, err
	}

	n += int64(TopicDataLength)

	o, err := w.Write(val)
	if err != nil {
		return n, err
	}

	n += int64(o)

	return n, nil
}

func (b *PubsubRequestBody) ToTLV() (tlvpac.TypeLengthValue, error) {
	raw := new(bytes.Buffer)
	_, err := b.WriteTo(raw)
	if err != nil {
		return nil, err
	}

	return tlvpac.TypeLengthValue(raw.Bytes()), nil
}

func (b *PubsubRequestBody) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := b.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}
