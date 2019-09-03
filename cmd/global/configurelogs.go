package global

import (
	"fmt"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type rootOptions struct {
	CfgFile  string
	flags    *pflag.FlagSet
	Debug    bool
	Hosts    []string
	LogLevel string

	SSHPrivateKey string
	DaemonAddress string
}

var Opt rootOptions

// ConfigureDaemonLogs sets the logrus logging level and formatting
func ConfigureDaemonLogs(Opt *rootOptions) error {
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
