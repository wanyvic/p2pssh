package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func PingHandle(conn *net.TCPConn, config api.ClientConfig) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			conn.Write([]byte("exit\n"))
		}
	}()
	r := io.Reader(conn)
	w := io.Writer(conn)

	b := []byte(fmt.Sprintf(p2p.P2PINGCONNECT+"\n%s\n", config.NodeID))
	conn.Write(b)

	go io.Copy(w, os.Stdin)
	io.Copy(os.Stdout, r)

	logrus.Info("handle exit")
}
