package api

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"golang.org/x/crypto/ssh"
)

type ClientConfig struct {
	UserName string
	NodeID   peer.ID
	Auth     []ssh.AuthMethod
}
