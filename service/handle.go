package service

import (
	"io"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func Handle(tcpConn *net.TCPConn) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	var buf [20480]byte
	n, err := reader.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	header := string(buf[:n])
	logrus.Debug(header)
	if strings.Contains(header, p2p.P2PSSHCONNECT) {
		SSHandle(tcpConn, header)
	} else if strings.Contains(header, p2p.P2PINGCONNECT) {
		PingHandle(tcpConn, header)
	}
	logrus.Debug("Handle Close")
}
