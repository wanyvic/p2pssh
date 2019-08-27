package cmd

import (
	"fmt"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	p2p "github.com/wanyvic/p2pssh/libp2p"
)

// configureDaemonLogs sets the logrus logging level and formatting
func configureDaemonLogs(Opt *daemonOptions) error {
	if Opt.LogLevel != "" {
		lvl, err := logrus.ParseLevel(Opt.LogLevel)
		if err != nil {
			return fmt.Errorf("unable to parse logging level: %s", Opt.LogLevel)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	p2p.PrivateKey = Opt.PrivateKey
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		DisableColors:   false,
		FullTimestamp:   true,
	})
	return nil
}
