package ssh

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/wanyvic/p2pssh/api"
	"golang.org/x/crypto/ssh"
)

func connect(userName string, password string, privateBytes []byte, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	if len(privateBytes) > 0 {
		Signer, err := ssh.ParsePrivateKey(privateBytes)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(Signer))
	}
	clientConfig = &ssh.ClientConfig{
		User:    userName,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
func Start(r io.Reader, w io.Writer, config api.ClientConfig) error {
	session, err := connect(config.UserName, config.Password, config.PrivateKey, "127.0.0.1", 22)
	if err != nil {
		return err
	}
	defer session.Close()
	// excute command
	session.Stdout = w
	session.Stderr = w
	session.Stdin = r

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", config.Height, config.Width, modes); err != nil {
		return errors.New("request for pseudo terminal failed: " + err.Error())
	}

	session.Run("bash")
	return nil
}
