package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	tlvpac "bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

const (
	KeyDataLength    uint8 = 16
	ValLenDataLength uint8 = 4
)

type RedisRequestBody struct {
	Key   string
	Len   uint32
	Value tlvpac.TypeLengthValue
}

func (b *RedisRequestBody) GetLen() uint32 {
	return uint32(KeyDataLength) + uint32(ValLenDataLength) + b.Len
}

func (b *RedisRequestBody) ReadFromIO(r io.Reader) error {
	var (
		typ  uint8
		blen uint32
		key  = make([]byte, 16)
		vlen uint32
		err  error
	)

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return err
	}

	if typ != tlvpac.RedisRequestPayloadType {
		return errors.New("Invalid redis request")
	}

	err = binary.Read(r, binary.BigEndian, &blen)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &key)
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

	*b = RedisRequestBody{
		Key:   string(key),
		Len:   vlen,
		Value: tlvpac.TypeLengthValue(buf),
	}

	return nil
}

func (b *RedisRequestBody) WriteToIO(w io.Writer) error {
	if len(b.Key) > 16 {
		return errors.New("Error key length exceed limit")
	}

	typ := tlvpac.RedisRequestPayloadType
	blen := b.GetLen()
	key := append(make([]byte, 16-len(b.Key)), []byte(b.Key)...)
	vlen := b.Len
	val := b.Value

	binary.Write(w, binary.BigEndian, typ)
	binary.Write(w, binary.BigEndian, blen)
	binary.Write(w, binary.BigEndian, []byte(key))
	binary.Write(w, binary.BigEndian, vlen)
	binary.Write(w, binary.BigEndian, val)

	return nil
}

func (b *RedisRequestBody) ToTLV() (tlvpac.TypeLengthValue, error) {
	raw := new(bytes.Buffer)
	b.WriteToIO(raw)

	return tlvpac.TypeLengthValue(raw.Bytes()), nil
}

func (b *RedisRequestBody) FromTLV(tlv tlvpac.TypeLengthValue) error {
	r := bytes.NewReader(tlv)
	err := b.ReadFromIO(r)
	if err != nil {
		return err
	}

	return nil
}
