package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"bitbucket.org/non-pn/mini-redis-go/cmd/client/internal/handler"
	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

const (
	cliGetCmd = "get"
	cliSetCmd = "set"
	cliSubCmd = "sub"
	cliPubCmd = "pub"
)

func handleClientCommand(cli *network.Client, cmd string, vals []string) error {
	if cli == nil {
		return errors.New("Error: client cannot be nil")
	}

	switch cmd {
	case cliGetCmd:
		if len(vals) == 0 {
			return errors.New("Error: Get cmd require at least one argument")
		}
		k := vals[0]
		handler.HandleClientGet(cli, k)
		break
	case cliSetCmd:
		if len(vals) < 2 {
			return errors.New("Error: Get cmd require two argument")
		}
		k, v := vals[0], vals[1]
		handler.HandleClientSet(cli, k, v)
		break
	case cliSubCmd:
		if len(vals) == 0 {
			return errors.New("Error: Sub cmd require at least one argument")
		}
		topic := vals[0]
		handler.HandleClientSub(cli, topic)
		break
	case cliPubCmd:
		if len(vals) < 2 {
			return errors.New("Error: Sub cmd require topic and message")
		}
		topic, message := vals[0], vals[1]
		handler.HandleClientPub(cli, topic, message)
		break
	default:
		fmt.Println("Invalid command")
		break
	}

	return nil
}

func main() {
	var (
		host string
		port string
		cert string
		key  string
		cmd  string
		vals []string
	)
	flag.StringVar(&host, "h", constant.DefaultServerHost, "host for client to connect to")
	flag.StringVar(&port, "p", constant.DefaultServerPort, "port for client to connect to")
	flag.StringVar(&cert, "cert", "", "absolute path to cert file for ssl")
	flag.StringVar(&key, "key", "", "absolute path to key file for ssl")
	flag.StringVar(&cmd, "c", "get", "command to use")

	flag.Parse()
	vals = flag.Args()

	// Init client and try connect to host
	port = ":" + port
	client := network.NewClient(constant.Protocol, host, port, cert, key)
	log.Println("Start redis client, try connecting to:", host+port)

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	handleClientCommand(client, cmd, vals)
}
