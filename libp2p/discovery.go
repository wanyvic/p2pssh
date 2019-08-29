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
		logrus.Warning("no connection,connecting bootstrap now")
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
					logrus.Error(err)
				} else {
					logrus.Debug("Connection established with bootstrap node:", *peerinfo)
				}
			}()
		}
		wg.Wait()
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
			logrus.Error(err)
			return
		}
		for pr := range peerChan {
			time.Sleep(time.Second)
			if pr.ID == p.host.ID() {
				continue
			}
			logrus.Debug("Found peer:", pr)
			if p.host.Network().Connectedness(pr.ID) != network.Connected {
				logrus.Debug("Connecting to:", pr)
				if err := p.host.Connect(context.Background(), pr); err != nil {
					logrus.Debug("Connection failed:", err)
					continue
				}
				logrus.Debug("Connection established peer:", pr)
			}
		}
	}
}
