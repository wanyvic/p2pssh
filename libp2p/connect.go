package p2p

import (
	"bufio"
	"context"
	"io"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	swarm "github.com/libp2p/go-libp2p-swarm"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func (p *P2PSSH) Connect(pid peer.ID, reader *bufio.Reader, writer *bufio.Writer) error {
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
	r := bufio.NewReader(stream)
	w := bufio.NewWriter(stream)

	go swap(w, reader)
	go swap(writer, r)
	return nil
}
func swap(dst *bufio.Writer, src *bufio.Reader) {
	var buf [1024]byte
	for {
		n, err := src.Read(buf[:])
		if err != nil || err == io.EOF {
			break
		}
		logrus.Debug(string(buf[:n]))
		dst.Write(buf[:n])
		dst.Flush()
	}
}
