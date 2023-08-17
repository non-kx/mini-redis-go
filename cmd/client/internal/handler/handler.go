package handler

import (
	"fmt"

	"bitbucket.org/non-pn/mini-redis-go/internal/network"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func HandleClientGet(cli *network.Client, k string) {
	resp, err := cli.Get(k)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", resp.String())
}

func HandleClientSet(cli *network.Client, k string, v string) {
	s := tlv.String(v)
	resp, err := cli.Set(k, &s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", resp)
}

func HandleClientSub(cli *network.Client, topic string) {
	sub, err := cli.Sub(topic)
	if err != nil {
		panic(err)
	}

	fmt.Println("Subscribed to topic:", topic)
	err = sub.Subscribe(func(s string) {
		fmt.Println("incoming msg:", s)
	})
	if err != nil {
		panic(err)
	}
}

func HandleClientPub(cli *network.Client, topic string, msg string) {
	resp, err := cli.Pub(topic, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
