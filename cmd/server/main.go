package main

import (
	"fmt"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
)

func main() {

	PORT := ":" + "6377"
	err := redis.InitRedisServer(PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Start redis server at port", PORT)

	err = redis.StartRedisServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer redis.StopRedisServer()
}
