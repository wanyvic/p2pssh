package service

import (
	"io"
	"net"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func PingHandle(tcpConn *net.TCPConn) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	writer := io.Writer(tcpConn)
	var buf [1024]byte
	libp2p := p2p.GetLibp2p()
	n, err := reader.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if nodeID, found := getNodeID(string(buf[:n])); found {
		err := libp2p.Ping(nodeID, reader, writer)
		if err != nil {
			logrus.Error("Connect to ", nodeID, " failed")
			return
		}
	}
	logrus.Debug("SSHandle Close")
}
func getNodeID(str string) (nodeID peer.ID, found bool) {
	var err error
	if strings.Contains(str, p2p.P2PINGCONNECT) {
		array := strings.Split(str, "\n")
		if len(array) < 3 {
			logrus.Error("no nodeID")
			return nodeID, false
		}
		nodeID, err = peer.IDB58Decode(array[1])
		if err != nil {
			return nodeID, false
		}
	}
	logrus.Debug("Unmarshal Peer ID ", nodeID)
	return nodeID, true
}
