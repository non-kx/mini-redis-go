package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func handleNewConn(conn net.Conn) {
	fmt.Println("New connection")
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(string(netData)) == "Ping" {
			conn.Write([]byte("Pong\n"))
			continue
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
}

func run() {
	fmt.Println("Hello from server")

	PORT := ":" + "6377"
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel() //Avoid context leak

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleNewConn(c)
	}
}
