package redis

import (
	"errors"
	"fmt"
	"strconv"

	"bitbucket.org/non-pn/mini-redis-go/internal/network"
)

type RedisClient struct {
	netcli *network.Client
}

func NewClient(host string) *RedisClient {
	rediscli := &RedisClient{
		netcli: network.NewClient(PROTOCOL, host),
	}

	return rediscli
}

func (cli *RedisClient) Connect() error {
	err := cli.netcli.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (cli *RedisClient) Close() error {
	err := cli.netcli.Close()
	if err != nil {
		return err
	}

	return nil
}

func (cli *RedisClient) SendGetCmd(k string) ([]byte, error) {
	pkg, err := createPayload(GetCmd, &k, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.netcli.Send(string(pkg))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func (cli *RedisClient) SendSetCmd(k string, v string) ([]byte, error) {
	pkg, err := createPayload(SetCmd, &k, &v)
	if err != nil {
		return nil, err
	}

	resp, err := cli.netcli.Send(string(pkg))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createPayload(cmd RedisCmd, k *string, v *string) ([]byte, error) {
	var (
		cmdbyte []byte
		kbyte   []byte
		vbyte   []byte
	)
	cmdbyte = []byte(strconv.Itoa(int(cmd)))

	if k != nil {
		kbyte = []byte(*k)
		if len(kbyte) > 8 {
			return nil, errors.New("Invalid key length")
		}

		zbyte := make([]byte, 8-len(kbyte))
		kbyte = append(zbyte, kbyte...)
	}

	if v != nil {
		vbyte = []byte(*v)
	}

	return append(append(cmdbyte, kbyte...), vbyte...), nil
}
