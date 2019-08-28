package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

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

type client struct {
	Addr   *net.TCPAddr
	conn   *net.TCPConn
	config api.ClientConfig
	ctx    context.Context
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
	c.handle()
	return nil
}
func (c *client) Close() {
	c.conn.Close()
}
func (c *client) handle() {

	reader := bufio.NewReader(c.conn)
	writer := bufio.NewWriter(c.conn)
	auth, _ := json.Marshal(c.config)
	b := []byte(fmt.Sprintf("--------P2PSSH--CONNECT--------\n%s\n", string(auth)))
	c.conn.Write(b)
	time.Sleep(time.Second)
	c.conn.Write(b)
	go io.Copy(writer, os.Stdin)
	go io.Copy(os.Stdout, reader)

	logrus.Debug("exit")
}
