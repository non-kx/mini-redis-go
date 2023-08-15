package test

import (
	"fmt"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func StartServerForTest() {
	serv, err := network.NewServer(constant.PROTOCOL, constant.REDIS_SERVER_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = serv.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer serv.Stop()
}

func Test_receive_message_subscribed_topic(t *testing.T) {
	var (
		topic     = "test"
		msg       = "test_msg"
		subclient *network.Client
		pubclient *network.Client

		err error
	)

	subclient = network.NewClient(constant.PROTOCOL, constant.DEFAULT_REDIS_SEVER_HOST)
	err = subclient.Connect()
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer subclient.Close()

	pubclient = network.NewClient(constant.PROTOCOL, constant.DEFAULT_REDIS_SEVER_HOST)
	err = pubclient.Connect()
	if err != nil {
		t.Errorf("Got err = %v", err)
	}
	defer pubclient.Close()

	t.Run("it should receive published message", func(t *testing.T) {
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
}
