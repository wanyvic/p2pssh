package p2p

import "github.com/libp2p/go-libp2p/p2p/protocol/ping"

const (
	P2PINGCONNECT = "--------P2PING--CONNECT--------"
)

func (p *P2PSSH) NewPingService() {
	ping.NewPingService(p.host)
}
