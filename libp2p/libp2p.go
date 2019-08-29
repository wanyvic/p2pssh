package p2p

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
	"github.com/wanyvic/p2pssh/ssh"
)

var (
	bootstarp = []string{
		"/ip4/119.3.66.159/udp/9000/quic/p2p/QmdQERFyHXZE4mBUuSrjbcuicRrmrQk4BB6uTAfiFWWjvq",
		"/ip4/132.232.79.195/udp/9000/quic/p2p/QmZLdPPkXanNCaQYk7CUaQkPioYBnoanhHC4Z9ZvF7eNWt",
	}
)

const (
	MAX_CONNECTION = 8
	Community      = "p2pssh"
	ID             = "/ssh/1.0.0"
)

type P2PSSH struct {
	host             host.Host
	routingDiscovery *discovery.RoutingDiscovery
	dht              *kaddht.IpfsDHT
}

var (
	Libp2p     *P2PSSH
	PrivateKey string
)

func GetLibp2p() *P2PSSH {
	if Libp2p == nil {
		var err error
		Libp2p, err = New()
		if err != nil {
			logrus.Error(err)
			Libp2p = nil
		}
		return Libp2p
	}
	return Libp2p
}
func New() (_ *P2PSSH, err error) {
	p := P2PSSH{}

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
	)
	listenAddrs := libp2p.ListenAddrStrings(
		// "/ip4/0.0.0.0/tcp/9000",
		"/ip4/0.0.0.0/udp/9000/quic",
	)

	if PrivateKey != "" {
		priv_bytes, err := hex.DecodeString(PrivateKey)
		if err != nil {
			return nil, err
		}
		priv, err := crypto.UnmarshalPrivateKey(priv_bytes)
		if err != nil {
			return nil, err
		}
		p.host, err = libp2p.New(context.Background(), transports, listenAddrs, libp2p.Identity(priv), libp2p.EnableRelay(circuit.OptDiscovery), libp2p.NATPortMap())
	} else {
		p.host, err = libp2p.New(context.Background(), transports, listenAddrs, libp2p.EnableRelay(circuit.OptDiscovery), libp2p.NATPortMap())
	}
	if err != nil {
		return nil, err
	}
	logrus.Info(p.host.ID(), p.host.Addrs())

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
	go p.connectBootstarp()
	go p.connectFromDHT()
	p.NewSSHService()
	return &p, nil
}
func (p *P2PSSH) connectBootstarp() {
	var wg sync.WaitGroup
	for _, peerAddr := range bootstarp {
		maddr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			fmt.Println(err)
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
	if len(p.host.Network().Peers()) < 1 {
		logrus.Warning("no connection")
	}
}
func (p *P2PSSH) connectFromDHT() {
	for len(p.host.Network().Peers()) < MAX_CONNECTION {
		peerChan, err := p.routingDiscovery.FindPeers(context.Background(), Community)
		if err != nil {
			panic(err)
		}
		for pr := range peerChan {
			time.Sleep(time.Second)
			if pr.ID == p.host.ID() {
				continue
			}
			fmt.Println("Found peer:", pr)
			if p.host.Network().Connectedness(pr.ID) != network.Connected {
				fmt.Println("Connecting to:", pr)
				if err := p.host.Connect(context.Background(), pr); err != nil {
					fmt.Println("Connection failed:", err)
					continue
				}
				fmt.Println("Connection established peer:", pr)
			}
		}
		break
	}
}
func (p *P2PSSH) NewSSHService() {
	p.host.SetStreamHandler(ID, handleStream)
}

func handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	r := io.Reader(s)
	w := io.Writer(s)
	var buf [1024]byte

	n, err := r.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug(string(buf[:n]))
	if auth, found := parse(string(buf[:n])); found {
		if err := ssh.Start(r, w, auth); err != nil {
			w.Write([]byte(err.Error()))
		}
	}
	logrus.Debug("stream close")
	s.Close()

	// stream 's' will stay open until you close it (or the other side closes it).
}
func parse(str string) (auth api.ClientConfig, found bool) {
	if strings.Contains(str, "--------P2PSSH--CONNECT--------") {
		logrus.Debug("finding peer id")
		svar := strings.Split(str, "\n")
		if len(svar) < 3 {
			logrus.Error("no auth")
			return auth, false
		}
		json.Unmarshal([]byte(svar[1]), &auth)
	}
	return auth, true
}
