package service

import (
	"net"
)

const (
	DefaultListenAddr = "127.0.0.1:9001"
)

var (
	DefaultConnect = func() *net.TCPAddr {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", DefaultListenAddr)
		return tcpAddr
	}
)

type ConnHandler func(*net.TCPConn)

type p2pService struct {
	TCPAddr     *net.TCPAddr
	ConnHandler ConnHandler
}

type Service interface {
	Start() error
	SetHandler(handle ConnHandler)
	Close()
}
