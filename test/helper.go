package test

import (
	"testing"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func startTestServer(t *testing.T) {
	serv, err := network.NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "", "")
	if err != nil {
		t.Error(err)
	}
	go func(t *testing.T) {
		err = serv.Start()
		if err != nil {
			t.Error(err)
		}
		defer serv.Stop()
	}(t)

	time.Sleep(5000)
}

func startTestSecureServer(t *testing.T) {
	serv, err := network.NewServer(constant.Protocol, ":"+constant.DefaultServerPort, "", "")
	if err != nil {
		t.Error(err)
	}
	go func(t *testing.T) {
		err = serv.Start()
		if err != nil {
			t.Error(err)
		}
		defer serv.Stop()
	}(t)

	time.Sleep(5000)
}

func createTestClient(t *testing.T) *network.Client {
	client := network.NewClient(constant.Protocol, constant.DefaultServerUrl)
	err := client.Connect("", "")
	if err != nil {
		t.Errorf("Got err = %v", err)
	}

	return client
}
