package p2p

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	"github.com/wanyvic/p2pssh/ssh"
)

const (
	P2PSSHCONNECT   = "--------P2PSSH--CONNECT--------"
	P2PSSHCONNECTED = "--------P2PSSH--CONNECTED--------"
	ID              = "/ssh/1.0.0"
)

func (p *P2PSSH) NewSSHService() {
	p.host.SetStreamHandler(ID, handleStream)
}

func handleStream(s network.Stream) {
	logrus.Debug("Got a new ssh stream!")
	// Create a buffer stream for non blocking read and write.
	r := io.Reader(s)
	w := io.Writer(s)
	scanner := bufio.NewScanner(s)
	if scanner.Scan() {
		_ = scanner.Text()
		if auth, found := UnmarshalConfig(scanner); found {
			if err := ssh.Start(r, w, auth); err != nil {
				logrus.Error(err)
				w.Write([]byte(err.Error()))
			}
		}
	}
	logrus.Debug("ssh stream close")
	s.Close()

	// stream 's' will stay open until you close it (or the other side closes it).
}
func UnmarshalConfig(scanner *bufio.Scanner) (auth api.ClientConfig, found bool) {
	if scanner.Scan() {
		str := scanner.Text()
		logrus.Debug("receive <-- ", str)
		err := json.Unmarshal([]byte(str), &auth)
		if err != nil {
			logrus.Error(err)
		}
		return auth, true
	}
	return auth, false
}
