package global

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

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
		TimestampFormat: RFC3339NanoFixed,
		DisableColors:   false,
		FullTimestamp:   true,
	})
	return nil
}
