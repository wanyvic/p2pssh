package rpcserver

import (
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func (s *Server) ConnectLS(req api.ConnectRequests, res *api.ConnectResponses) error {
	host := p2p.GetLibp2p()
	res.Peers = host.GetPeers()
	logrus.Debug("send --> ", res)
	return nil
}

func (s *Server) ConnectAdd(req api.ConnectAddRequests, res *api.ConnectAddResponses) error {
	host := p2p.GetLibp2p()
	err := host.ConnectPeerInfo(&req.PeerAddr)
	if err != nil {
		logrus.Error(err)
		res.Err = "Connection failed"
	} else {
		res.Result = "Connected successfully"
	}
	logrus.Debug("send --> ", res)
	return nil
}

func (s *Server) ConnectRm(req api.ConnectRmRequests, res *api.ConnectRmResponses) error {
	host := p2p.GetLibp2p()
	err := host.DisConnectPeer(req.Peer)
	if err != nil {
		logrus.Error(err)
		res.Err = "Connection Disconnect failed"
	} else {
		res.Result = "Connection Disconnect successfully"
	}
	logrus.Debug("send --> ", res)
	return nil
}
