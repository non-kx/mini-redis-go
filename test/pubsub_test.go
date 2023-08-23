package test

import (
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func Test_receive_message_subscribed_topic(t *testing.T) {
	var (
		topic     = "test"
		msg       = "test_msg"
		subclient *network.Client
		pubclient *network.Client

		err error
	)

	subclient = network.NewClient(constant.Protocol, constant.DefaultServerHost, ":"+constant.DefaultServerPort, "", "")
	err = subclient.Connect()
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer subclient.Close()

	pubclient = network.NewClient(constant.Protocol, constant.DefaultServerHost, ":"+constant.DefaultServerPort, "", "")
	err = pubclient.Connect()
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer pubclient.Close()

	t.Run("t should receive published message", func(t *testing.T) {
		subscriber, err := subclient.Sub(topic)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		resp, err := pubclient.Pub(topic, msg)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		if resp != "OK" {
			t.Errorf("Expect want = %v got = %v", "OK", resp)
		}

		rec, err := subscriber.NextMessage()
		if err != nil {
			t.Errorf("Async() err = %v", err)
		}

		if resp != "OK" {
			t.Errorf("Expect want = %v got = %v", msg, rec)
		}
	})
	t.Run("it should receive published message for multiple subscriber", func(t *testing.T) {
		c1 := createTestClient(t)
		c2 := createTestClient(t)
		defer c1.Close()
		defer c2.Close()

		sub1, err := c1.Sub(topic)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}
		sub2, err := c2.Sub(topic)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		resp, err := pubclient.Pub(topic, msg)
		if err != nil {
			t.Errorf("Got err = %v", err)
		}

		if resp != "OK" {
			t.Errorf("Expect want = %v got = %v", "OK", resp)
		}

		rec1, err := sub1.NextMessage()
		if err != nil {
			t.Errorf("Async() err = %v", err)
		}
		if rec1.String() != msg {
			t.Errorf("Expect want = %v got = %v", msg, rec1)
		}

		rec2, err := sub2.NextMessage()
		if err != nil {
			t.Errorf("Async() err = %v", err)
		}
		if rec2.String() != msg {
			t.Errorf("Expect want = %v got = %v", msg, rec2)
		}
	})

}
