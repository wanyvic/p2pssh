package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func New(ctx context.Context, tcpAddr *net.TCPAddr) p2pService {
	logrus.Debug("New")
	s := p2pService{}
	s.TCPAddr = tcpAddr
	s.ctx = ctx
	return s
}
func (s *p2pService) Start() error {
	if s.ConnHandler == nil {
		return errors.New("ConnHandler not set")
	}
	listener, err := net.ListenTCP("tcp", s.TCPAddr)
	defer listener.Close()
	if err != nil {
		return err
	}
	for {
		select {
		case <-time.After(time.Second):
			break
		case <-s.ctx.Done():
			return nil
		}
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Info("A client connected :" + tcpConn.RemoteAddr().String())
		go s.ConnHandler(tcpConn)
	}
}
func (s *p2pService) Close() {
}

func (s *p2pService) SetHandler(handle ConnHandler) {
	s.ConnHandler = handle
}

func Handle(tcpConn *net.TCPConn) {
	reader := bufio.NewReader(tcpConn)
	writer := bufio.NewWriter(tcpConn)
	var buf [1024]byte
	libp2p := p2p.GetLibp2p()

	n, err := reader.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if auth, found := parse(string(buf[:n])); found {

		err := libp2p.Connect(auth.NodeID, reader, writer)
		if err != nil {
			return
		}
	}
	logrus.Debug("exit")
}

func parse(str string) (auth api.ClientConfig, found bool) {
	if strings.Contains(str, "--------P2PSSH--CONNECT--------") {
		logrus.Debug("finding peer id")
		svar := strings.Split(str, "\n")
		if len(svar) < 3 {
			logrus.Error("no auth")
			return auth, false
		}
		err := json.Unmarshal([]byte(svar[1]), &auth)
		if err != nil {
			logrus.Error(err)
		}
	}
	return auth, true
}
