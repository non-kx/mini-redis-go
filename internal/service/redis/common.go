package redis

type RedisCmd uint8

const (
	PROTOCOL = "tcp"
)

const (
	GetCmd RedisCmd = iota
	SetCmd
)

type RedisPayload struct {
	Cmd   RedisCmd
	Key   string
	Value string
}
