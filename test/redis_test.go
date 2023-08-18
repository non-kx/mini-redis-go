package test

import (
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
	"github.com/stretchr/testify/assert"
)

func Test_redis(t *testing.T) {
	var (
		client *network.Client

		err error
	)

	client = network.NewClient(constant.Protocol, constant.DefaultServerUrl)
	err = client.Connect("", "")
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer client.Close()

	t.Run("it should return empty value for unset key", func(t *testing.T) {
		k := "unset"
		val, err := client.Get(k)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		assert.Equal(t, "", val.String())
	})
	t.Run("it should response OK for set command", func(t *testing.T) {
		k := "test"
		v := tlv.String("test")
		resp, err := client.Set(k, &v)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		assert.Equal(t, "OK", resp)
	})
	t.Run("it should set and get correctly", func(t *testing.T) {
		k := "test2"
		v := tlv.String("test2")
		resp, err := client.Set(k, &v)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		assert.Equal(t, "OK", resp)

		val, err := client.Get(k)
		assert.Equal(t, v.String(), val.String())
	})
	// t.Run("it should not be able to use key length more that 16 character", func(t *testing.T) {
	// 	k := "0123456789abcdefg"
	// 	v := tlv.String("some")
	// 	_, err := client.Set(k, &v)
	// 	if assert.Error(t, err) {
	// 		assert.EqualError(t, err, "")
	// 	}
	// })
}
