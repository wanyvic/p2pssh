package service

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func SSHandle(tcpConn *net.TCPConn) {
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
	writer.Write([]byte(p2p.P2PSSHCONNECTED))
	if auth, found := p2p.UnmarshalConfig(string(buf[:n])); found {
		err := libp2p.Connect(auth.NodeID, reader, writer)
		if err != nil {
			logrus.Error("Connect to ", auth.NodeID, " failed")
			return
		}
	}
	logrus.Debug("SSHandle Close")
}
