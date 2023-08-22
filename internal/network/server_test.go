package network

import (
	"path/filepath"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	serv, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "", "")
	defer serv.Stop()

	assert.Nil(t, err)
	assert.NotNil(t, serv)
	assert.NotNil(t, serv.RedisDb)
	assert.NotNil(t, serv.PubsubDb)
	assert.Equal(t, 0, len(serv.Connections))
}

func TestNewSecureServerWithInvalidCert(t *testing.T) {
	_, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "invalid", "invalid")

	assert.NotNil(t, err)
}

func TestNewSecureServerWithValidCert(t *testing.T) {
	cert, err := filepath.Abs("../../test/cert/server/server.crt")
	key, err := filepath.Abs("../../test/cert/server/server.key")

	assert.Nil(t, err)

	serv, err := NewServer(constant.Protocol, ":"+constant.DefaultServerPort, cert, key)
	defer serv.Stop()

	assert.Nil(t, err)
	assert.NotNil(t, serv)
	assert.NotNil(t, serv.RedisDb)
	assert.NotNil(t, serv.PubsubDb)
	assert.Equal(t, 0, len(serv.Connections))
}

func TestHandleConnection(t *testing.T) {

}
