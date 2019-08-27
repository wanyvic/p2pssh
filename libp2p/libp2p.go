package p2p

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"os"
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
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

var (
	bootstarp = []string{
		"/ip4/119.3.66.159/tcp/9000/p2p/QmdQERFyHXZE4mBUuSrjbcuicRrmrQk4BB6uTAfiFWWjvq",
		"/ip4/132.232.79.195/tcp/9000/p2p/QmZLdPPkXanNCaQYk7CUaQkPioYBnoanhHC4Z9ZvF7eNWt",
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
		Libp2p, err := New()
		if err != nil {
			logrus.Error(err)
			Libp2p = nil
		}
		return Libp2p
	}
	return Libp2p
}
func New() (p *P2PSSH, err error) {
	p = &P2PSSH{}
	if PrivateKey != "" {
		priv_bytes, err := hex.DecodeString(PrivateKey)
		if err != nil {
			return nil, err
		}
		priv, err := crypto.UnmarshalPrivateKey(priv_bytes)
		if err != nil {
			return nil, err
		}
		p.host, err = libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"), libp2p.Identity(priv), libp2p.EnableRelay(circuit.OptDiscovery), libp2p.NATPortMap())
	} else {
		p.host, err = libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"), libp2p.EnableRelay(circuit.OptDiscovery), libp2p.NATPortMap())
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	logrus.Info(p.host.ID(), p.host.Addrs())

	p.dht, err = kaddht.New(context.Background(), p.host)
	if err != nil {
		fmt.Println(err)
	}
	err = p.dht.Bootstrap(context.Background())
	if err != nil {
		panic(err)
	}

	p.routingDiscovery = discovery.NewRoutingDiscovery(p.dht)
	discovery.Advertise(context.Background(), p.routingDiscovery, Community)
	p.connectBootstarp()
	go p.connectFromDHT()
	p.NewSSHService()
	return p, nil
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
	r := bufio.NewReader(s)
	w := bufio.NewWriter(s)
	go readData(r)
	go writeData(w)
	// ssh.Start(r, w)

	// stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.Reader) {
	for {

		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.Writer) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		sendData, _, err := stdReader.ReadLine()

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}
