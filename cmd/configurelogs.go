package cmd

import (
	"fmt"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
)

// configureDaemonLogs sets the logrus logging level and formatting
func configureDaemonLogs(Opt *rootOptions) error {
	if Opt.LogLevel != "" {
		lvl, err := logrus.ParseLevel(Opt.LogLevel)
		if err != nil {
			return fmt.Errorf("unable to parse logging level: %s", Opt.LogLevel)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		DisableColors:   false,
		FullTimestamp:   true,
	})
	return nil
}
