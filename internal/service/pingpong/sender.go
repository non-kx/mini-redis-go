package pingpong

import (
	"net"

	"bitbucket.org/non-pn/mini-redis-go/internal/payload"
	"bitbucket.org/non-pn/mini-redis-go/internal/tools/tlv"
)

func SendPingRequest(conn net.Conn, msg *string) (*tlv.String, error) {
	if msg == nil {
		ping := string("PING")
		msg = &ping
	}
	s := tlv.String(*msg)
	msgtlv, err := s.ToTLV()
	if err != nil {
		return nil, err
	}

	req := payload.RequestPayload{
		Cmd:  payload.PingCmd,
		Body: msgtlv,
	}

	_, err = req.WriteTo(conn)
	if err != nil {
		return nil, err
	}

	resp := tlv.String("")
	_, err = resp.ReadFrom(conn)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
