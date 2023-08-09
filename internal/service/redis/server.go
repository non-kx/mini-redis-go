package redis

import (
	"log"
	"strconv"

	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

var (
	RedisServer *network.Server
	RedisDB     *db.KVStore
)

func StartRedisServer() error {
	err := RedisServer.Listen()
	if err != nil {
		return err
	}

	return nil
}

func StopRedisServer() error {
	err := RedisServer.Close()
	if err != nil {
		return err
	}

	return nil
}

func Get(k string) []byte {
	log.Printf("Get from KVStore with k[%v]", k)
	return RedisDB.Get(k)
}

func Set(k string, v []byte) {
	log.Printf("Set to KVStore with k[%v] v[%v]\n", k, v)
	RedisDB.Set(k, v)
}

func RedisHandler(ctx *network.RequestContext) error {
	data := ctx.Data
	pkg, err := transformPackage(data)
	if err != nil {
		return err
	}

	switch pkg.Cmd {
	case GetCmd:
		err = ctx.Response(Get(pkg.Key))
		if err != nil {
			log.Println(err)
			return err
		}
		break
	case SetCmd:
		Set(pkg.Key, []byte(pkg.Value))
		err = ctx.Response([]byte("Success"))
		if err != nil {
			log.Println(err)
			return err
		}
		break
	}

	err = RedisDB.CacheStorage()
	if err != nil {
		log.Println(err)
	}

	return nil
}

func InitRedisServer(port string) error {
	server, err := network.NewServer("tcp", port, RedisHandler)
	if err != nil {
		return err
	}

	kvstore := db.InitKVStore(nil)
	RedisServer = server
	RedisDB = kvstore

	return nil
}

func transformPackage(data []byte) (*RedisPayload, error) {
	command, err := strconv.Atoi(string(data[:1]))
	if err != nil {
		return nil, err
	}

	key := string(data[1:9])
	value := string(data[9:])
	return &RedisPayload{
		Cmd:   RedisCmd(command),
		Key:   key,
		Value: value,
	}, nil
}

// type RedisServer struct {
// 	network.Server
// 	DB *db.KVStore
// }

// func NewRedisServer(port string) error {
// 	kvstore := db.InitKVStore(nil)
// 	l, err := net.Listen(PROTOCOL, port)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &RedisServer{
// 		Server{
// 			Port:        port,
// 			Listener:    l,
// 			Connections: make([]*net.Conn, 0, 5),
// 			Handler:     handler,
// 		},
// 	}, nil

// 	return nil
// }
