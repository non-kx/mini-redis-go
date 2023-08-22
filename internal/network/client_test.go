package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("tcp", "localhost")

	assert.NotNil(t, client)
	assert.Equal(t, "tcp", client.Network)
	assert.Equal(t, "localhost", client.Host)
	assert.Nil(t, client.Connection)
}
