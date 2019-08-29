package client

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
)

const (
	DefaultConnectAddr = "127.0.0.1:9001"
)

var (
	DefaultConnect = func() *net.TCPAddr {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", DefaultConnectAddr)
		return tcpAddr
	}
)

type ConnHandler func(*net.TCPConn, api.ClientConfig)
type client struct {
	Addr        *net.TCPAddr
	conn        *net.TCPConn
	config      api.ClientConfig
	ctx         context.Context
	ConnHandler ConnHandler
}

func New(ctx context.Context, tcpAddr *net.TCPAddr, config api.ClientConfig) (c client) {
	logrus.Debug("Create New Client")
	c.Addr = tcpAddr
	c.config = config
	c.ctx = ctx
	return c
}
func (c *client) Connect() error {
	logrus.Debug("Client Connectting ...")
	conn, err := net.DialTCP("tcp", nil, c.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.ConnHandler(conn, c.config)
	return nil
}
func (c *client) Close() {
	c.conn.Close()
}
