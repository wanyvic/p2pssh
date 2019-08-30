package p2p

import (
	"encoding/json"
	"io"
	"strings"

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
	var buf [20480]byte

	n, err := r.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if auth, found := UnmarshalConfig(string(buf[:n])); found {
		if err := ssh.Start(r, w, auth); err != nil {
			logrus.Error(err)
			w.Write([]byte(err.Error()))
		}
	}
	logrus.Debug("ssh stream close")
	s.Close()

	// stream 's' will stay open until you close it (or the other side closes it).
}
func UnmarshalConfig(str string) (auth api.ClientConfig, found bool) {
	if strings.Contains(str, P2PSSHCONNECT) {
		array := strings.Split(str, "\n")
		if len(array) < 2 {
			logrus.Error("no auth")
			return auth, false
		}
		err := json.Unmarshal([]byte(array[1]), &auth)
		if err != nil {
			logrus.Error(err)
		}
	}
	logrus.Debug("Unmarshal Peer ID ", auth.NodeID)
	return auth, true
}
