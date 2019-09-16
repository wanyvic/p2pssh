package p2p

import (
	"bufio"
	"encoding/json"

	"github.com/sirupsen/logrus"
	lssh "github.com/wanyvic/go-libp2p-ssh"
	"github.com/wanyvic/p2pssh/api"
)

const (
	P2PSSHCONNECT   = "--------P2PSSH--CONNECT--------"
	P2PSSHCONNECTED = "--------P2PSSH--CONNECTED--------"
	ID              = "/ssh/1.0.0"
)

func (p *P2PSSH) NewSSHService() {
	config, err := lssh.DefaultServerConfig()
	if err != nil {
		logrus.Error(err)
	}
	lssh.NewSSHService(p.host, config)
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
