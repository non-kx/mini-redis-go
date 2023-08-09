package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
)

const (
	DEFAULT_REDIS_SEVER_HOST = "127.0.0.1:6377"
)

func main() {

	client := redis.NewClient(DEFAULT_REDIS_SEVER_HOST)
	log.Println("Start redis client, try connecting to host", DEFAULT_REDIS_SEVER_HOST)

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	for {
		var (
			resp   []byte
			reader = bufio.NewReader(os.Stdin)
		)
		fmt.Print(">> ")
		line, _ := reader.ReadString('\n')
		cmd := strings.Split(strings.TrimSuffix(line, "\n"), " ")

		switch cmd[0] {
		case "get":
			k := cmd[1]
			log.Println("Getting from redis with key:", k)
			resp, err = client.SendGetCmd(k)
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Printf("Raw response from server: %v\n", resp)
		case "set":
			k := cmd[1]
			v := cmd[2]
			log.Println("Setting to redis with key and value:", k, v)
			resp, err = client.SendSetCmd(k, v)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Println("->: " + string(resp))
	}
}
