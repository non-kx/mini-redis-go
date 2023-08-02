package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

func main() {
	fmt.Println("Hello from client")

	client := network.NewClient("tcp", "127.0.0.1:6377")
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		resp, err := client.Send(text)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("->: " + string(resp))
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
