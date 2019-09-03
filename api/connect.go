package api

import peer "github.com/libp2p/go-libp2p-peer"

type ConnectLSRequests struct {
}
type ConnectLSResponses struct {
	Peers []peer.ID
}
