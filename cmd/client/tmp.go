package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
// 	"bitbucket.org/non-pn/mini-redis-go/internal/service/redis"
// 	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
// )

// func main() {
// 	client := redis.NewClient(constant.DEFAULT_REDIS_SEVER_HOST)
// 	log.Println("Start redis client, try connecting to host", constant.DEFAULT_REDIS_SEVER_HOST)

// 	err := client.Connect()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer client.Close()

// 	for {
// 		var (
// 			resp   *redis.RedisResponsePayload
// 			reader = bufio.NewReader(os.Stdin)
// 		)
// 		fmt.Print(">> ")
// 		line, _ := reader.ReadString('\n')
// 		cmd := strings.Split(strings.TrimSuffix(line, "\n"), " ")

// 		switch cmd[0] {
// 		case "get":
// 			k := cmd[1]
// 			log.Println("Getting from redis with key:", k)

// 			resp, err = client.SendGetCmd(k)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			log.Printf("Raw response from server: %v\n", resp)

// 			tlvresp := tlv.TypeLengthValue(resp.RespBody)
// 			switch tlvresp.GetType() {
// 			case tlv.BinaryType:
// 				b := tlv.Binary([]byte{})
// 				err := b.FromTLV(tlvresp)
// 				if err != nil {
// 					fmt.Println(err)
// 					return
// 				}
// 				fmt.Printf("->: %v\n", b)
// 			case tlv.StringType:
// 				s := tlv.String("")
// 				err := s.FromTLV(tlvresp)
// 				if err != nil {
// 					fmt.Println(err)
// 					return
// 				}
// 				fmt.Printf("->: %v\n", s)
// 			default:
// 				fmt.Printf("->: %v\n", "nil")
// 				break
// 			}
// 		case "set":
// 			k := cmd[1]
// 			v := tools.String(string(cmd[2]))
// 			log.Println("Setting to redis with key and value:", k, v)

// 			tlv, err := v.ToTLV()
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			resp, err = client.SendSetCmd(k, tlv)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			fmt.Printf("->: %v\n", string(resp.RespBody))
// 		}
// 	}
// }
