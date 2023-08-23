package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTLVType(t *testing.T) {
	stlv := TypeLengthValue([]byte{2, 0, 0, 0, 4, 116, 101, 115, 116})
	typ := stlv.GetType()

	assert.Equal(t, StringType, typ)
}

func TestGetEmptyTLVType(t *testing.T) {
	stlv := TypeLengthValue([]byte{})
	typ := stlv.GetType()

	assert.Equal(t, EmptyType, typ)
}

func TestGetTLVLength(t *testing.T) {
	stlv := TypeLengthValue([]byte{2, 0, 0, 0, 4, 116, 101, 115, 116})
	slen := stlv.GetLength()

	assert.Equal(t, 4, int(slen))
}

func TestGetEmptyTLVLength(t *testing.T) {
	stlv := TypeLengthValue([]byte{})
	slen := stlv.GetLength()

	assert.Equal(t, 0, int(slen))
}

func TestGetTLVValue(t *testing.T) {
	stlv := TypeLengthValue([]byte{2, 0, 0, 0, 4, 116, 101, 115, 116})
	sval := stlv.GetValue()

	assert.Equal(t, []byte(stlv)[5:], sval)
}

func TestGetEmptyTLVValue(t *testing.T) {
	stlv := TypeLengthValue([]byte{})
	sval := stlv.GetValue()

	assert.Equal(t, []byte{}, sval)
}
