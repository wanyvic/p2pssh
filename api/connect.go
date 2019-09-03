package api

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

type ConnectRequests struct {
}
type ConnectResponses struct {
	Peers []peer.ID
}
type ConnectAddRequests struct {
	PeerAddr peer.AddrInfo
}
type ConnectAddResponses struct {
	Result string
	Err    string
}

type ConnectRmRequests struct {
	Peer peer.ID
}
type ConnectRmResponses ConnectAddResponses
