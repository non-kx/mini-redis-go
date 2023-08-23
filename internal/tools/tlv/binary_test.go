package tlv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryReadFrom(t *testing.T) {
	testb := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	testreader := bytes.NewReader(testb)

	b := new(Binary)
	n, err := b.ReadFrom(testreader)

	assert.Equal(t, len(testb), int(n))
	assert.Nil(t, err)
	assert.Equal(t, testb[5:], []byte(*b))
}

func TestBinaryReadFromInvalid(t *testing.T) {
	testb := []byte{2, 0, 0, 0, 4, 116, 101, 115, 116}
	testreader := bytes.NewReader(testb)

	b := new(Binary)
	n, err := b.ReadFrom(testreader)

	assert.Equal(t, 1, int(n))
	assert.NotNil(t, err)
}

func TestBinaryWriteTo(t *testing.T) {
	testb := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	testwriter := new(bytes.Buffer)

	b := Binary(testb[5:])
	n, err := b.WriteTo(testwriter)

	buf := testwriter.Bytes()

	assert.Equal(t, len(testb), int(n))
	assert.Nil(t, err)
	assert.Equal(t, buf, testb)
}

func TestBinaryFromTLV(t *testing.T) {
	tlv := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	b := new(Binary)

	err := b.FromTLV(tlv)

	assert.Nil(t, err)
	assert.Equal(t, tlv[5:], []byte(*b))
}

func TestBinaryToTLV(t *testing.T) {
	testtlv := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	b := Binary(testtlv[5:])

	tlv, err := b.ToTLV()

	assert.Nil(t, err)
	assert.Equal(t, testtlv, []byte(tlv))
}

func TestBinaryToString(t *testing.T) {
	testtlv := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	b := Binary(testtlv[5:])
	s := b.String()

	assert.Equal(t, "test", s)
}
