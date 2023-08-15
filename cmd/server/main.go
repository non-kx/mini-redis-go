package main

import (
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func main() {
	serv, err := network.NewServer(constant.PROTOCOL, constant.REDIS_SERVER_PORT)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Start redis server at port", constant.REDIS_SERVER_PORT)
	err = serv.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer serv.Stop()
}
