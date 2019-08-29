package service

import (
	"fmt"
	"io"
	"net"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func PingHandle(tcpConn *net.TCPConn, header string) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	writer := io.Writer(tcpConn)
	libp2p := p2p.GetLibp2p()

	if nodeID, found := getNodeID(header); found {
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
	fmt.Println("i")
	if strings.Contains(str, p2p.P2PINGCONNECT) {
		fmt.Println("in")
		array := strings.Split(str, "\n")
		if len(array) < 3 {
			logrus.Error("no nodeID")
			return nodeID, false
		}
		fmt.Println(array[0])
		fmt.Println(array[1])
		fmt.Println(array[2])
		nodeID, err = peer.IDB58Decode(array[1])
		if err != nil {
			return nodeID, false
		}
	}
	logrus.Debug("Unmarshal Peer ID ", nodeID)
	return nodeID, true
}
