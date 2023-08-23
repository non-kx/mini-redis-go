package tlv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringReadFrom(t *testing.T) {
	tests := []byte{2, 0, 0, 0, 4, 116, 101, 115, 116}
	testreader := bytes.NewReader(tests)

	s := new(String)
	n, err := s.ReadFrom(testreader)

	assert.Equal(t, len(tests), int(n))
	assert.Nil(t, err)
	assert.Equal(t, tests[5:], []byte(*s))
}

func TestStringReadFromInvalid(t *testing.T) {
	tests := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116}
	testreader := bytes.NewReader(tests)

	s := new(String)
	n, err := s.ReadFrom(testreader)

	assert.Equal(t, 1, int(n))
	assert.NotNil(t, err)
}

func TestStringWriteTo(t *testing.T) {
	tests := []byte{2, 0, 0, 0, 4, 116, 101, 115, 116}
	testwriter := new(bytes.Buffer)

	s := String(tests[5:])
	n, err := s.WriteTo(testwriter)

	buf := testwriter.Bytes()

	assert.Equal(t, len(tests), int(n))
	assert.Nil(t, err)
	assert.Equal(t, buf, tests)
}

func TestStringFromTLV(t *testing.T) {
	tlv := []byte{2, 0, 0, 0, 4, 116, 101, 115, 116}
	s := new(String)

	err := s.FromTLV(tlv)

	assert.Nil(t, err)
	assert.Equal(t, string(tlv[5:]), s.String())
}

func TestStringToTLV(t *testing.T) {
	testtlv := []byte{2, 0, 0, 0, 4, 116, 101, 115, 116}
	s := String(testtlv[5:])

	tlv, err := s.ToTLV()

	assert.Nil(t, err)
	assert.Equal(t, testtlv, []byte(tlv))
}

func TestStringToString(t *testing.T) {
	s := String("test_string")
	ss := s.String()

	assert.Equal(t, "test_string", ss)
}

func TestNilStringToString(t *testing.T) {
	s := new(String)
	ss := s.String()

	assert.Equal(t, "", ss)
}
