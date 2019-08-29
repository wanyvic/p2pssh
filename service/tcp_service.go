package service

import (
	"context"
	"errors"
	"net"

	"github.com/sirupsen/logrus"
)

func New(ctx context.Context, tcpAddr *net.TCPAddr) p2pService {
	logrus.Debug("New")
	s := p2pService{}
	s.TCPAddr = tcpAddr
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
