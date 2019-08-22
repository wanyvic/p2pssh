package network

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	crypto "github.com/libp2p/go-libp2p-crypto"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	swarm "github.com/libp2p/go-libp2p-swarm"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/wanyvic/prizes_ssh/ssh"
)

var (
	bootstarp = []string{
		"/ip4/119.3.66.159/tcp/9000/p2p/QmdQERFyHXZE4mBUuSrjbcuicRrmrQk4BB6uTAfiFWWjvq",
		"/ip4/132.232.79.195/tcp/9000/p2p/QmZLdPPkXanNCaQYk7CUaQkPioYBnoanhHC4Z9ZvF7eNWt",
	}
)

const (
	MAX_CONNECTION = 8
	Community      = "prizes_ssh"
	ID             = "/ssh/1.0.0"
)

type P2PSSH struct {
	host             host.Host
	routingDiscovery *discovery.RoutingDiscovery
	dht              *kaddht.IpfsDHT
}

func New(priva crypto.PrivKey) (_ *P2PSSH, err error) {
	p := P2PSSH{}

	p.host, err = libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9001"), libp2p.Identity(priva), libp2p.EnableRelay(circuit.OptDiscovery), libp2p.NATPortMap())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(p.host.ID(), p.host.Addrs())

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
				fmt.Println(err)
			} else {
				fmt.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()
	if len(p.host.Network().Peers()) < 1 {
		fmt.Println("no connection")
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
func (p *P2PSSH) Connect(pid peer.ID) error {
	fmt.Println("connect")
	addrInfo, err := p.dht.FindPeer(context.Background(), pid)
	if err != nil {
		fmt.Println(pid)
		return err
	}
	if err := p.host.Connect(context.Background(), addrInfo); err != nil {
		fmt.Println("Connection failed:", err)
		p.host.Network().(*swarm.Swarm).Backoff().Clear(addrInfo.ID)
		relayaddr, _ := ma.NewMultiaddr("/p2p-circuit/ipfs/" + addrInfo.ID.Pretty())
		h3relayInfo := peer.AddrInfo{
			ID:    addrInfo.ID,
			Addrs: []ma.Multiaddr{relayaddr},
		}
		if err := p.host.Connect(context.Background(), h3relayInfo); err != nil {
			fmt.Println(h3relayInfo, err)
		}
	}
	fmt.Println("Connection established peer:", addrInfo)

	stream, err := p.host.NewStream(context.Background(), addrInfo.ID, protocol.ID(ID))

	if err != nil {
		fmt.Println("Connection failed:", err)
		return err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go writeData(rw)
	go readData(rw)
	return nil
}

func (p *P2PSSH) NewSSHService(h host.Host) {
	h.SetStreamHandler(ID, handleStream)
}

func handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	ssh.Start(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
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

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}
