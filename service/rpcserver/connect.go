package rpcserver

import (
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func (s *Server) ConnectLS(req api.ConnectLSRequests, res *api.ConnectLSResponses) error {
	host := p2p.GetLibp2p()
	res.Peers = host.GetPeers()
	logrus.Debug("send --> ", res.Peers)
	return nil
}
