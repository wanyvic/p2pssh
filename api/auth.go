package api

import peer "github.com/libp2p/go-libp2p-peer"

type UserAuth struct {
	User     string
	Password string
	Pubkey   string
	NodeID   peer.ID
}
