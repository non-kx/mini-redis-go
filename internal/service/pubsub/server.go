package pubsub

import (
	"io"
	"log"
	"net"
	"time"

	"bitbucket.org/non-pn/mini-redis-go/internal/constant"
	"bitbucket.org/non-pn/mini-redis-go/internal/db"
	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

type PubsubServer struct {
	Port        string
	Listener    net.Listener
	Connections []*net.Conn
	DB          *db.KVStore[[]*net.Conn]
}

func NewServer(port string) (*PubsubServer, error) {
	kvstore := db.InitKVStore[[]*net.Conn](nil)
	l, err := net.Listen(constant.PROTOCOL, port)
	if err != nil {
		return nil, err
	}

	return &PubsubServer{
		Port:        port,
		Listener:    l,
		Connections: make([]*net.Conn, 0, 5),
		DB:          kvstore,
	}, nil
}

func (serv *PubsubServer) Start() error {
	for {
		c, err := serv.Listener.Accept()
		if err != nil {
			return err
		}

		log.Println("New connection from:", c.RemoteAddr())
		go serv.HandleConnection(&c)
	}
}

func (serv *PubsubServer) Stop() error {
	for _, conn := range serv.Connections {
		err := (*conn).Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (serv *PubsubServer) HandleConnection(conn *net.Conn) error {
	serv.Connections = append(serv.Connections, conn)

	for {
		data, err := network.ReadUntilCRLF(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			}
			return err
		}

		reqctx := &network.RequestContext{
			Now:  time.Now(),
			Data: data,
			Conn: conn,
		}
		err = serv.HandleRequest(reqctx)
		if err != nil {
			return err
		}
	}
}

func (serv *PubsubServer) HandleRequest(ctx *network.RequestContext) error {
	var resp PubsubResponsePayload
	data := PubsubRawRequestPayload(ctx.Data)
	payload, err := data.TransformPayload()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Receive payload:", payload)

	switch payload.Cmd {
	case SubscribeCmd:
		resp = PubsubResponsePayload{
			RespType: SubscribeSuccess,
			RespBody: []byte("OK"),
		}
		raw, err := resp.ToRaw()
		if err != nil {
			log.Println(err)
			return err
		}
		err = ctx.Response(*raw)
		if err != nil {
			log.Println(err)
			return err
		}
		break
	case PublishCmd:
		// serv.Set(payload.Key, []byte(payload.Value))
		// resp = RedisResponsePayload{
		// 	RespType: SetSuccess,
		// 	RespBody: []byte("OK"),
		// }
		// raw, err := resp.ToRaw()
		// if err != nil {
		// 	log.Println(err)
		// 	return err
		// }
		// err = ctx.Response(*raw)
		// if err != nil {
		// 	log.Println(err)
		// 	return err
		// }
		// break
	}

	err = serv.DB.CacheStorage()
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (serv *PubsubServer) SubscribeToTopic(topic string, conn *net.Conn) {
	conns := serv.DB.Get(topic)
	conns = append(conns, conn)
	serv.DB.Set(topic, conns)
}

func (serv *PubsubServer) PublishToTopic(topic string, msg []byte) {

}
