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
	Len   uint32
	Value tlvpac.TypeLengthValue
}

func (b *PubsubRequestBody) GetLen() uint32 {
	return uint32(TopicDataLength) + uint32(ValLenDataLength) + b.Len
}

func (b *PubsubRequestBody) ReadFromIO(r io.Reader) error {
	var (
		typ   uint8
		blen  uint32
		topic = make([]byte, 16)
		vlen  uint32
		err   error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return err
	}

	if typ != tlvpac.RedisRequestPayloadType {
		return errors.New("Invalid pubsub request")
	}

	err = binary.Read(r, binary.BigEndian, &blen)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &topic)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &vlen)
	if err != nil {
		return err
	}

	buf := make([]byte, vlen)

	if vlen > 0 {
		_, err = r.Read(buf)
		if err != nil {
			return err
		}
	}

	*b = PubsubRequestBody{
		Topic: string(topic),
		Len:   vlen,
		Value: tlvpac.TypeLengthValue(buf),
	}

	return nil
}

func (b *PubsubRequestBody) WriteToIO(w io.Writer) error {
	if len(b.Topic) > 16 {
		return errors.New("Error topic length exceed limit")
	}

	typ := tlvpac.RedisRequestPayloadType
	blen := b.GetLen()
	topic := append(make([]byte, 16-len(b.Topic)), []byte(b.Topic)...)
	vlen := b.Len
	val := b.Value

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, blen)
	binary.Write(w, binary.BigEndian, []byte(topic))
	binary.Write(w, binary.BigEndian, vlen)
	binary.Write(w, binary.BigEndian, val)

	return nil
}

func (b *PubsubRequestBody) ToTLV() (tlvpac.TypeLengthValue, error) {
	raw := new(bytes.Buffer)
	b.WriteToIO(raw)

	return tlvpac.TypeLengthValue(raw.Bytes()), nil
}

func (b *PubsubRequestBody) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	err := b.ReadFromIO(r)
	if err != nil {
		return err
	}

	return nil
}
