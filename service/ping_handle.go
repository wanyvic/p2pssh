package service

import (
	"bufio"
	"io"
	"net"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func PingHandle(tcpConn *net.TCPConn, nodeID peer.ID) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	writer := io.Writer(tcpConn)
	libp2p := p2p.GetLibp2p()

	err := libp2p.Ping(nodeID, reader, writer)
	if err != nil {
		logrus.Error("Connect to ", nodeID, " failed")
		writer.Write([]byte("Connect failed"))
		return
	}

	logrus.Debug("PingHandle Close")
}
func getNodeID(scanner *bufio.Scanner) (nodeID peer.ID, found bool) {
	if scanner.Scan() {
		str := scanner.Text()
		logrus.Debug(str)
		nodeID, err := peer.IDB58Decode(str)
		if err != nil {
			return nodeID, false
		}
		logrus.Debug("Unmarshal Peer ID ", nodeID)
		return nodeID, true
	}
	return nodeID, false
}
