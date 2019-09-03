package client

import (
	"io"
	"net"
	"net/rpc/jsonrpc"
	"strings"

	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

func JsonRPConnect(address string, Method string, args interface{}, reply interface{}) error {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	conn.Write([]byte(p2p.P2PRPCONNECT + "\n"))
	buf := make([]byte, 1024)
	r := io.Reader(conn)
	n, err := r.Read(buf[:])
	if err != nil || err == io.EOF {
		logrus.Error(err)
	}
	logrus.Debug("receive <-- ", string(buf[:n]))
	if strings.Contains(string(buf[:n]), p2p.P2PRPCONNECTED) {
		jsonConn := jsonrpc.NewClient(conn)
		err = jsonConn.Call(Method, args, reply)
		if err != nil {
			return err
		}
	}

	return nil
}
