package service

import (
	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

type RedisServer struct {
	Server *network.Server
	DB     *db.KVStore
}

func (rd *RedisServer) Start() {
	rd.Server.Listen()
}

func NewRedisServer(port string) (*RedisServer, error) {
	server, err := network.NewServer("tcp", port)
	if err != nil {
		return nil, err
	}

	kvstore, err := db.InitKVStore(nil)
	if err != nil {
		return nil, err
	}

	return &RedisServer{
		Server: server,
		DB:     kvstore,
	}, nil
}
