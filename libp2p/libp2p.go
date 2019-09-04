package p2p

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/p2p/login"
)

const (
	Community = "p2pssh"
)

type P2PSSH struct {
	host             host.Host
	routingDiscovery *discovery.RoutingDiscovery
	dht              *kaddht.IpfsDHT
	bootstrap        []string
}

var (
	Libp2p     *P2PSSH
	PrivateKey string
)

func GetLibp2p() *P2PSSH {
	if Libp2p == nil {
		var err error
		Libp2p, err = New(PrivateKey)
		if err != nil {
			logrus.Error(err)
			Libp2p = nil
		}
		return Libp2p
	}
	return Libp2p
}
func New(privkey string) (_ *P2PSSH, err error) {
	p := P2PSSH{}
	transports, listenAddrs := login.GetTransport()

	priv, err := getPrivateKey(privkey)
	if err != nil {
		return nil, err
	}
	p.host, err = libp2p.New(context.Background(), transports, listenAddrs, libp2p.Identity(priv), libp2p.EnableRelay(circuit.OptDiscovery, circuit.OptHop, circuit.OptActive), libp2p.NATPortMap())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Your PeerID is :%s\nListen:%s\n", p.host.ID().String(), p.host.Addrs())

	p.dht, err = kaddht.New(context.Background(), p.host)
	if err != nil {
		return nil, err
	}
	err = p.dht.Bootstrap(context.Background())
	if err != nil {
		return nil, err
	}

	p.routingDiscovery = discovery.NewRoutingDiscovery(p.dht)
	discovery.Advertise(context.Background(), p.routingDiscovery, Community)
	p.bootstrap = BootStrap
	go p.ConnectionManager()
	return &p, nil
}

func getPrivateKey(privkey string) (priv crypto.PrivKey, err error) {
	if privkey == "" {
		priv, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	} else {
		privBytes, err := hex.DecodeString(privkey)
		if err == nil {
			priv, err = crypto.UnmarshalPrivateKey(privBytes)
		}
	}
	if err != nil {
		return nil, err
	}
	return priv, nil
}
func (p *P2PSSH) GetPeers() []peer.ID {
	return p.host.Network().Peers()
}
