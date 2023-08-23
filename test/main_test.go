package test

import (
	"os"
	"testing"

	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

var serv *network.Server

// func init() {
// 	secure := flag.Bool("secure", false, "Enable secure to tcp connection")
// 	flag.Parse()

// 	if *secure {
// 		os.Setenv("SECURE", "1")
// 	}
// }

func beforeAll() {
	serv = createTestServer()
}

func afterAll() {
	serv.Stop()
}

func TestMain(m *testing.M) {
	beforeAll()
	exitVal := m.Run()
	// afterAll()

	os.Exit(exitVal)
}
