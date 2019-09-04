package login

import (
	"errors"
	"os"
	"os/exec"

	"github.com/libp2p/go-libp2p"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	"golang.org/x/crypto/ssh/terminal"
)

func getTerminalSize() (int, int, error) {
	fd := int(0)
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		fd = int(os.Stdin.Fd())
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return 0, 0, errors.New(err.Error() + "error allocating terminal")
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return 0, 0, err
	}
	defer terminal.Restore(fd, oldState)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return 0, 0, err
	}
	return termWidth, termHeight, nil
}
func SetTerminalEcho(flag bool) {
	if flag {

		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	} else {
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	}
}
func GetTransport() (libp2p.Option, libp2p.Option) {
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
	)
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/9000",
		"/ip4/0.0.0.0/udp/9000/quic",
	)
	return transports, listenAddrs
}
