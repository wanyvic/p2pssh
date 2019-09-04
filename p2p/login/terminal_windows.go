package login

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-tcp-transport"
	"github.com/sirupsen/logrus"
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
	var out bytes.Buffer
	cmd := exec.Command("mode", "con")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: error reading console width %s", err.Error())
	}
	re := regexp.MustCompile(`\d+`)
	rs := re.FindAllString(out.String(), -1)
	i, err := strconv.Atoi(rs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: error transfering string to int %s", err.Error())
	}
	j, err := strconv.Atoi(rs[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: error transfering string to int %s", err.Error())
	}
	logrus.Debug(rs, "width: ", j, " line: ", i)
	return j, i, nil
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
