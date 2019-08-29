package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
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
func Ping(tcpAddr *net.TCPAddr, NodeID peer.ID) error {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	_ = conn
	return nil
}
func (c *client) handle() {

	r := io.Reader(c.conn)
	w := io.Writer(c.conn)
	auth, err := json.Marshal(c.config)
	if err != nil {
		logrus.Error(err)
	}
	b := []byte(fmt.Sprintf("--------P2PSSH--CONNECT--------\n%s\n", string(auth)))
	time.Sleep(time.Second)
	c.conn.Write(b)
	var buf [1024]byte

	n, err := r.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if strings.Contains(string(buf[:n]), "--------P2PSSH--CONNECTED--------") {
		c.conn.Write(b)
		go io.Copy(w, os.Stdin)
		io.Copy(os.Stdout, r)
	}

}
