package p2p

import (
	"context"
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	MAX_CONNECTION       = 8
	BOOTSTRAP_CONNECTION = 2
)

func (p *P2PSSH) ConnectionManager() {
	for {
		if len(p.host.Network().Peers()) <= 0 {
			p.connectFromBootstarp()
		} else {
			p.connectFromDHT()
		}
		time.Sleep(time.Second * 3)
	}
}
func (p *P2PSSH) getRandomFromBootstrap() int64 {
	bigInt := big.NewInt(int64(len(p.bootstrap)))
	randomInt, _ := rand.Int(rand.Reader, bigInt)
	return randomInt.Int64()
}
func (p *P2PSSH) connectFromBootstarp() {
	var wg sync.WaitGroup
	for len(p.host.Network().Peers()) <= 0 {
		logrus.Debug("no connection,connecting bootstrap now")
		for i := 0; i < BOOTSTRAP_CONNECTION; i++ {
			index := p.getRandomFromBootstrap()
			maddr, err := ma.NewMultiaddr(p.bootstrap[index])
			if err != nil {
				logrus.Error(err)
				continue
			}
			peerinfo, _ := peer.AddrInfoFromP2pAddr(maddr)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := p.host.Connect(context.Background(), *peerinfo); err != nil {
					logrus.Error("Connection bootstrap failed: ", peerinfo.ID)
				} else {
					logrus.Info("Connection established with bootstrap node:", *peerinfo)
				}
			}()
		}
		wg.Wait()
		time.Sleep(time.Second * 3)
	}
}
func (p *P2PSSH) connectFromDHT() {
	for len(p.host.Network().Peers()) < MAX_CONNECTION {
		if len(p.host.Network().Peers()) <= 0 {
			logrus.Warning("no connection ,connect bootstrap first")
			return
		}
		peerChan, err := p.routingDiscovery.FindPeers(context.Background(), Community)
		if err != nil {
			logrus.Error("no found peers")
			return
		}
		for pr := range peerChan {
			time.Sleep(time.Second)
			if pr.ID == p.host.ID() {
				continue
			}
			if p.host.Network().Connectedness(pr.ID) != network.Connected {
				logrus.Debug("Connecting to:", pr)
				if err := p.host.Connect(context.Background(), pr); err != nil {
					logrus.Debug("Connection failed:", pr)
					continue
				}
				logrus.Info("Connection established peer:", pr)
			}
		}
	}
}
func (p *P2PSSH) ConnectPeerInfo(peerInfo *peer.AddrInfo) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	if p.host.Network().Connectedness(peerInfo.ID) != network.Connected {
		logrus.Debug("Connecting to:", peerInfo)
		if err := p.host.Connect(ctx, *peerInfo); err != nil {
			logrus.Debug("Connection failed:", peerInfo)
			return err
		}
		logrus.Info("Connection established peer:", peerInfo)
	}
	return nil
}
func (p *P2PSSH) DisConnectPeer(peerID peer.ID) error {
	if p.host.Network().Connectedness(peerID) == network.Connected {
		logrus.Debug("DisConnecting to:", peerID)
		if err := p.host.Network().ClosePeer(peerID); err != nil {
			logrus.Debug("Connection failed:", peerID)
			return err
		}
		logrus.Info("Connection DisConnected peer:", peerID)
	}
	return nil
}
