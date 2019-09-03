package rpcserver

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	p2p "github.com/wanyvic/p2pssh/libp2p"
)

type Server struct {
}

func init() {
	rpc.Register(new(Server))
}
func Handle(tcpConn *net.TCPConn) {

	tcpConn.Write([]byte(p2p.P2PRPCONNECTED + "\n"))
	jsonrpc.ServeConn(tcpConn)
}
