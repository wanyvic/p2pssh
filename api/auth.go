package api

import (
	peer "github.com/libp2p/go-libp2p-peer"
)

type ClientConfig struct {
	UserName   string
	NodeID     peer.ID
	Password   string
	PrivateKey []byte
}
