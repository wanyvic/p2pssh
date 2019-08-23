package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

const DefaultHost = "localhost"
const DefaultPort = 9001

type Server struct {
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
}

func NewServer(host string, port int) (server *Server, err error) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	server = &Server{tcpAddr: addr}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Stop()
	}()
	return server, nil
}

func (u *Server) Start() error {

	listener, err := net.ListenTCP("tcp", u.tcpAddr)
	if err != nil {
		return err
	}
	u.tcpListener = listener
	go handler(u.tcpListener)
	return nil
}

func (u *Server) Stop() {
	u.tcpListener.Close()
}

func handler(tcpListener *net.TCPListener) {
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}
	logrus.Info("tcp handler exit")
}
func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	i := 0
	for {
		var buf [1024]byte
		message, err := reader.Read(buf[:])
		if err != nil || err == io.EOF {
			break
		}
		logrus.Debug(string(message))

		b := []byte(msg)

		conn.Write(b)

		i++

		if i > 10 {
			break
		}
	}
}
