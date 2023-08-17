package main

import (
	"flag"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func main() {
	var (
		port string
		cert string
		key  string
	)

	flag.StringVar(&port, "p", constant.DefaultServerPort, "port that server will listen on")
	flag.StringVar(&cert, "cert", constant.DefaultServerPort, "absolute path to cert file for ssl")
	flag.StringVar(&key, "key", constant.DefaultServerPort, "absolute path to key file for ssl")
	flag.Parse()

	serv, err := network.NewServer(constant.Protocol, ":"+port, nil, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Start redis server at port", port)
	err = serv.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer serv.Stop()
}
