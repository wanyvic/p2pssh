package service

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func SSHandle(tcpConn *net.TCPConn, auth api.ClientConfig) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	writer := io.Writer(tcpConn)
	libp2p := p2p.GetLibp2p()
	writer.Write([]byte(p2p.P2PSSHCONNECTED))
	err := libp2p.Connect(auth.NodeID, reader, writer)
	if err != nil {
		logrus.Error("Connect to ", auth.NodeID, " failed")
		writer.Write([]byte("Connect failed"))
		return
	}
	logrus.Debug("SSHandle Close")
}
