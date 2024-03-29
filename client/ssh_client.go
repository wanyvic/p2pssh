package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
	"github.com/wanyvic/p2pssh/p2p/login"
)

func SSHandle(conn *net.TCPConn, config api.ClientConfig) {
	login.SetTerminalEcho(true)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			conn.Write([]byte{3})
		}
	}()
	r := io.Reader(conn)
	w := io.Writer(conn)
	auth, err := json.Marshal(config)
	if err != nil {
		logrus.Error(err)
	}
	header := []byte(fmt.Sprintf(p2p.P2PSSHCONNECT+"\n%s\n", string(auth)))
	logrus.Debug("send --> ", string(header))
	conn.Write(header)
	go io.Copy(w, os.Stdin)
	convertColor(r)
	login.SetTerminalEcho(false)
	logrus.Debug("SSHandle exit")
}
