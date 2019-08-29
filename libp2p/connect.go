package p2p

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	swarm "github.com/libp2p/go-libp2p-swarm"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func (p *P2PSSH) Connect(pid peer.ID, reader io.Reader, writer io.Writer) error {
	logrus.Debug("p2p connect")
	addrInfo, err := p.dht.FindPeer(context.Background(), pid)
	if err != nil {
		logrus.Error(pid)
		return err
	}
	if err := p.host.Connect(context.Background(), addrInfo); err != nil {
		logrus.Error("Connection failed:", err)
		p.host.Network().(*swarm.Swarm).Backoff().Clear(addrInfo.ID)
		relayaddr, _ := ma.NewMultiaddr("/p2p-circuit/ipfs/" + addrInfo.ID.Pretty())
		h3relayInfo := peer.AddrInfo{
			ID:    addrInfo.ID,
			Addrs: []ma.Multiaddr{relayaddr},
		}
		if err := p.host.Connect(context.Background(), h3relayInfo); err != nil {
			logrus.Error(h3relayInfo, err)
		}
	}
	logrus.Debug("Connection established peer:", addrInfo)

	stream, err := p.host.NewStream(context.Background(), addrInfo.ID, protocol.ID(ID))

	if err != nil {
		logrus.Error("Connection failed:", err)
		return err
	}
	r := io.Reader(stream)
	w := io.Writer(stream)
	go io.Copy(w, reader)
	io.Copy(writer, r)
	logrus.Debug("libp2p ssh close")
	return nil
}
func (p *P2PSSH) Ping(pid peer.ID, reader io.Reader, writer io.Writer) error {
	logrus.Debug("p2p connect")
	addrInfo, err := p.dht.FindPeer(context.Background(), pid)
	if err != nil {
		logrus.Error(pid)
		return err
	}
	if err := p.host.Connect(context.Background(), addrInfo); err != nil {
		logrus.Warning("Connection failed:", err, " trying p2p-circuit/ipfs/")
		p.host.Network().(*swarm.Swarm).Backoff().Clear(addrInfo.ID)
		relayaddr, _ := ma.NewMultiaddr("/p2p-circuit/ipfs/" + addrInfo.ID.Pretty())
		h3relayInfo := peer.AddrInfo{
			ID:    addrInfo.ID,
			Addrs: []ma.Multiaddr{relayaddr},
		}
		if err := p.host.Connect(context.Background(), h3relayInfo); err != nil {
			logrus.Error(h3relayInfo, err)
		}
	}
	logrus.Debug("Connection established peer:", addrInfo)

	pctx, cancel := context.WithCancel(context.Background())
	go func() {
		var buf [1024]byte
		_, err := reader.Read(buf[:])
		if err != nil || err == io.EOF {
			logrus.Error(err)
			cancel()
		}
	}()
	defer cancel()
	ts := ping.Ping(pctx, p.host, pid)
	for {
		select {
		case res := <-ts:
			if res.Error != nil {
				return nil
			}
			writer.Write([]byte(fmt.Sprintf("ping took: %s\n", res.RTT.String())))
		case <-time.After(time.Second * 4):
			writer.Write([]byte("failed to receive ping\n"))
		case <-pctx.Done():
			logrus.Debug("libp2p ping close")
			return nil
		}
		time.Sleep(time.Second)
	}
}
