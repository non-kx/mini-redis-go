package tlv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorReadFrom(t *testing.T) {
	testerr := NewError(1, "test")
	testerrraw := []byte{5, 0, 0, 0, 6, 0, 1, 116, 101, 115, 116}
	testreader := bytes.NewReader(testerrraw)

	e := new(Error)
	n, err := e.ReadFrom(testreader)

	assert.Equal(t, len(testerrraw), int(n))
	assert.Nil(t, err)
	assert.Equal(t, testerr, *e)
}

func TestErrorReadFromInvalid(t *testing.T) {
	testerrraw := []byte{1, 0, 0, 0, 6, 0, 1, 116, 101, 115, 116}
	testreader := bytes.NewReader(testerrraw)

	e := new(Error)
	n, err := e.ReadFrom(testreader)

	assert.Equal(t, 1, int(n))
	assert.NotNil(t, err)
}

func TestErorrWriteTo(t *testing.T) {
	testerr := NewError(1, "test")
	testerrraw := []byte{5, 0, 0, 0, 6, 0, 1, 116, 101, 115, 116}
	testwriter := new(bytes.Buffer)

	n, err := testerr.WriteTo(testwriter)

	assert.Equal(t, len(testerrraw), int(n))
	assert.Nil(t, err)
	assert.Equal(t, testwriter.Bytes(), testerrraw)
}

func TestErrorFromTLV(t *testing.T) {
	tlv := []byte{5, 0, 0, 0, 6, 0, 1, 116, 101, 115, 116}
	e := new(Error)

	err := e.FromTLV(tlv)

	assert.Nil(t, err)
	assert.Equal(t, 1, int(e.Code))
	assert.Equal(t, "test", e.Msg)
}

func TestErrorToTLV(t *testing.T) {
	testtlv := []byte{5, 0, 0, 0, 6, 0, 1, 116, 101, 115, 116}
	e := Error{
		Code: 1,
		Msg:  "test",
	}

	tlv, err := e.ToTLV()

	assert.Nil(t, err)
	assert.Equal(t, testtlv, []byte(tlv))
}
