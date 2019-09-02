package login

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/howeyc/gopass"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/api"
)

func ParseClientConfig(str string, keyPath string) (config api.ClientConfig, err error) {
	if str == "" {
		return config, errors.New("no args")
	}
	nodeID := ""
	split := strings.Split(str, "@")
	if len(split) == 2 {
		config.UserName = split[0]
		nodeID = split[1]
	} else if len(split) == 1 {
		nodeID = split[0]
	} else {
		return config, errors.New("args errors")
	}
	config.NodeID, err = peer.IDB58Decode(nodeID)
	if err != nil {
		return config, err
	}
	if keyPath == "" {
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return config, err
		}
		logrus.Debug("Password: ", string(pass))
		config.Password = string(pass)
	} else {
		privateBytes, err := parsePrivateKey(keyPath)
		if err != nil {
			return config, err
		}
		config.PrivateKey = privateBytes
	}
	config.Width, config.Height, err = getTerminalSize()
	if err != nil {
		return config, err
	}
	return config, nil
}

func parsePrivateKey(keyPath string) (private []byte, _ error) {
	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return private, err
	}
	logrus.Debug("SSHPrivateKey: ", string(privateBytes))

	return privateBytes, nil
}
