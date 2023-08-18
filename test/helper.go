package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func secureMode() bool {
	secure := os.Getenv("SECURE")
	return secure == "1"
}

func getServerCert() (string, string) {
	if secureMode() {
		cert, err := filepath.Abs("../cert/server/server.crt")
		if err != nil {
			panic(err)
		}

		key, err := filepath.Abs("../cert/server/server.key")
		if err != nil {
			panic(err)
		}

		return cert, key
	}

	return "", ""
}

func getClientCert() (string, string) {
	if secureMode() {
		cert, err := filepath.Abs("../cert/client/client.crt")
		if err != nil {
			panic(err)
		}

		key, err := filepath.Abs("../cert/client/client.key")
		if err != nil {
			panic(err)
		}

		return cert, key
	}

	return "", ""
}

func createTestServer() *network.Server {
	cert, key := getServerCert()
	serv, err := network.NewServer(constant.Protocol, ":"+constant.DefaultServerPort, cert, key)
	if err != nil {
		panic(err)
	}
	go func() {
		err = serv.Start()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(5000)

	return serv
}

func createTestClient(t *testing.T) *network.Client {
	cert, key := getClientCert()
	client := network.NewClient(constant.Protocol, constant.DefaultServerUrl)
	err := client.Connect(cert, key)
	if err != nil {
		t.Errorf("Got err = %v", err)
	}

	return client
}
