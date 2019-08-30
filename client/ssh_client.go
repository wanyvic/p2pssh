package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func SSHandle(conn *net.TCPConn, config api.ClientConfig) {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
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
	b := []byte(fmt.Sprintf(p2p.P2PSSHCONNECT+"\n%s\n", string(auth)))
	time.Sleep(time.Second)
	conn.Write(b)
	buf := make([]byte, 1024)

	n, err := r.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if strings.Contains(string(buf[:n]), p2p.P2PSSHCONNECTED) {
		conn.Write(b)
		go io.Copy(w, os.Stdin)
		io.Copy(os.Stdout, r)
	}
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	logrus.Debug("SSHandle exit")
}
