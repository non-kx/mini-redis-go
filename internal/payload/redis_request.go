package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	tlvpac "bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

const (
	KeyDataLength uint8 = 16
)

type RedisRequestBody struct {
	Key   string
	Value tlvpac.TypeLengthValue
}

func (b *RedisRequestBody) ReadFrom(r io.Reader) (int64, error) {
	var (
		typ uint8
		len uint32
		key = make([]byte, KeyDataLength)
		n   int64

		err error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	n += 1

	if typ != tlvpac.RedisRequestPayloadType {
		return n, errors.New("Invalid redis request")
	}

	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return n, err
	}

	n += 4

	err = binary.Read(r, binary.BigEndian, &key)
	if err != nil {
		return n, err
	}

	n += int64(KeyDataLength)

	vlen := len - uint32(KeyDataLength)
	buf := make([]byte, vlen)

	if vlen > 0 {
		o, err := r.Read(buf)
		if err != nil {
			return n, err
		}

		n += int64(o)
	}

	*b = RedisRequestBody{
		Key:   string(key),
		Value: tlvpac.TypeLengthValue(buf),
	}

	return n, nil
}

func (b *RedisRequestBody) WriteTo(w io.Writer) (int64, error) {
	var (
		n int64

		err error
	)

	if len(b.Key) > int(KeyDataLength) {
		return 0, errors.New("Error: Key length exceed limit")
	}

	typ := tlvpac.RedisRequestPayloadType
	key := append(make([]byte, int(TopicDataLength)-len(b.Key)), []byte(b.Key)...)
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

	err = binary.Write(w, binary.BigEndian, key)
	if err != nil {
		return n, err
	}

	n += int64(KeyDataLength)

	o, err := w.Write(val)
	if err != nil {
		return n, err
	}

	n += int64(o)

	return n, nil
}

func (b *RedisRequestBody) ToTLV() (tlvpac.TypeLengthValue, error) {
	raw := new(bytes.Buffer)
	_, err := b.WriteTo(raw)
	if err != nil {
		return nil, err
	}

	return tlvpac.TypeLengthValue(raw.Bytes()), nil
}

func (b *RedisRequestBody) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	_, err := b.ReadFrom(r)
	if err != nil {
		return err
	}

	return nil
}
