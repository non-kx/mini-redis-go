package main

import (
	"fmt"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func main() {
	client := network.NewClient(constant.PROTOCOL, constant.DEFAULT_REDIS_SEVER_HOST)
	log.Println("Start redis client, try connecting to host", constant.DEFAULT_REDIS_SEVER_HOST)

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	// k := "PING"
	// // v := "PONG"
	// // s := tlv.String(v)
	// resp, err := client.Get(k)
	resp, err := client.Pub("test", "test_msg")
	if err != nil {
		fmt.Println(err)
	}

	// err = sub.Subscribe()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// msg := sub.NextMessage()

	fmt.Println("->:", resp)
}
