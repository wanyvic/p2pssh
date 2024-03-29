package login

import (
	"os"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-tcp-transport"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

var (
	dll            = syscall.MustLoadDLL("kernel32")
	setConsoleMode = dll.MustFindProc("SetConsoleMode")
	m              uint32
)

func init() {

	h := syscall.Handle(os.Stdin.Fd())
	if err := syscall.GetConsoleMode(h, &m); err != nil {
		logrus.Error(err)
	}
}
func getTerminalSize() (int, int, error) {
	if h, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE); err != nil {
		logrus.Error(err)
		return 0, 0, err
	} else {
		var info windows.ConsoleScreenBufferInfo
		if err := windows.GetConsoleScreenBufferInfo(h, &info); err != nil {
			logrus.Error(err)
			return 0, 0, err
		}
		width := info.Window.Right - info.Window.Left + 1
		height := info.Window.Bottom - info.Window.Top + 1
		return int(width), int(height), nil
	}
}
func SetTerminalEcho(flag bool) {
	h := syscall.Handle(os.Stdin.Fd())
	if flag {
		if err := SetInputConsoleMode(h, 0); err != nil {
			logrus.Error(err)
		}
	} else {
		if err := SetInputConsoleMode(h, m); err != nil {
			logrus.Error(err)
		}
	}
}

func SetInputConsoleMode(h syscall.Handle, m uint32) error {
	r, _, err := setConsoleMode.Call(uintptr(h), uintptr(m))
	if r == 0 {
		return err
	}
	return nil
}
func GetTransport() (libp2p.Option, libp2p.Option) {
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
	)
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/9000",
	)
	return transports, listenAddrs
}
