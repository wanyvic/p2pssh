package service

import (
	"io"
	"net"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
	"github.com/wanyvic/ssh"
)

func SSHandle(tcpConn *net.TCPConn, config api.ClientConfig) {
	defer tcpConn.Close()
	reader := io.Reader(tcpConn)
	writer := io.Writer(tcpConn)
	libp2p := p2p.GetLibp2p()

	auth := make([]ssh.AuthMethod, 0)
	if config.Password != "" {
		auth = append(auth, ssh.Password(config.Password))
	}
	if len(config.PrivateKey) > 0 {
		Signer, err := ssh.ParsePrivateKey(config.PrivateKey)
		if err != nil {
			logrus.Error(err)
			return
		}
		auth = append(auth, ssh.PublicKeys(Signer))
	}
	clientConfig := ssh.ClientConfig{
		User:    config.UserName,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote ma.Multiaddr, key ssh.PublicKey) error {
			return nil
		},
	}

	err := libp2p.SSHConnect(config.NodeID, clientConfig, reader, writer)
	if err != nil {
		logrus.Error("Connect to ", config.NodeID, " failed")
		writer.Write([]byte("Connect failed"))
		return
	}
	logrus.Debug("SSHandle Close")
}
