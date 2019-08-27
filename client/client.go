package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

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
	Addr     *net.TCPAddr
	conn     *net.TCPConn
	userAuth api.UserAuth
	ctx      context.Context
}

func New(ctx context.Context, tcpAddr *net.TCPAddr, userAuth api.UserAuth) (c client) {
	logrus.Debug("New")
	c.Addr = tcpAddr
	c.userAuth = userAuth
	c.ctx = ctx
	return c
}
func (c *client) Connect() error {
	logrus.Debug("Connect")
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
	auth, _ := json.Marshal(c.userAuth)
	b := []byte(fmt.Sprintf("--------P2PSSH--CONNECT--------\n%s\n", string(auth)))
	c.conn.Write(b)
	go io.Copy(writer, os.Stdin)
	go io.Copy(os.Stdout, reader)

	logrus.Debug("exit")
}
