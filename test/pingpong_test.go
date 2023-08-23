package test

import (
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
	"github.com/stretchr/testify/assert"
)

func Test_pingpong(t *testing.T) {
	var (
		client *network.Client

		err error
	)

	client = network.NewClient(constant.Protocol, constant.DefaultServerHost, ":"+constant.DefaultServerPort, "", "")
	err = client.Connect()
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer client.Close()

	t.Run("it should return PONG", func(t *testing.T) {
		s := "PING"
		resp, err := client.Ping(&s)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		assert.Equal(t, "PONG", resp)
	})
	t.Run("it should return same value", func(t *testing.T) {
		s := "test"
		resp, err := client.Ping(&s)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		assert.Equal(t, s, resp)
	})
}

func Test_secure_pingpong(t *testing.T) {

}
