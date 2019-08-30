package service

import (
	"bufio"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func Handle(tcpConn *net.TCPConn) {
	defer tcpConn.Close()
	scanner := bufio.NewScanner(tcpConn)
	if scanner.Scan() {
		str := scanner.Text()
		logrus.Debug("receive <-- ", str)
		if strings.Contains(str, p2p.P2PSSHCONNECT) {
			if auth, found := p2p.UnmarshalConfig(scanner); found {
				SSHandle(tcpConn, auth)
			}

		} else if strings.Contains(str, p2p.P2PINGCONNECT) {
			if nodeID, found := getNodeID(scanner); found {
				PingHandle(tcpConn, nodeID)
			}
		}
	}

	logrus.Debug("Handle Close")
}
