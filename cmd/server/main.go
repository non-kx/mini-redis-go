package main

import (
	"fmt"

	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func main() {
	fmt.Println("Hello from server")

	PORT := ":" + "6377"

	s, err := network.NewServer("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Listen()
}
