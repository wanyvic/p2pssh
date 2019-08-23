// +build !windows

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/cli/debug"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/wanyvic/dssh/cmd/config"
	"github.com/wanyvic/dssh/dssh/dsshd"
)

type DaemonCli struct {
	*config.Config
	configFile *string
	flags      *pflag.FlagSet
}

// NewDaemonCli returns a daemon CLI
func NewDaemonCli() *DaemonCli {
	return &DaemonCli{}
}

func getDefaultDaemonConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return "", errors.New("could not get either XDG_CONFIG_HOME or HOME")
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "prizes"), nil
}

func getDefaultDaemonConfigFile() (string, error) {
	dir, err := getDefaultDaemonConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "daemon.json"), nil
}

func getDaemonConfDir(_ string) (string, error) {
	return getDefaultDaemonConfigDir()
}

// installConfigFlags adds flags to the pflag.FlagSet to configure the daemon
func installConfigFlags(conf *config.Config, flags *pflag.FlagSet) error {
	return nil
}

func runDaemon(opts *daemonOptions) error {
	daemonCli := NewDaemonCli()
	return daemonCli.start(opts)
}

func initLogging(_, stderr io.Writer) {
	logrus.SetOutput(stderr)
}
func loadDaemonCliConfig(opts *daemonOptions) (*config.Config, error) {
	conf := opts.daemonConfig
	conf.LogLevel = opts.LogLevel
	return conf, nil
}
func (cli *DaemonCli) start(opts *daemonOptions) (err error) {
	stopc := make(chan bool)
	defer close(stopc)

	opts.SetDefaultOptions(opts.flags)

	if cli.Config, err = loadDaemonCliConfig(opts); err != nil {
		return err
	}

	if err := configureDaemonLogs(cli.Config); err != nil {
		return err
	}

	cli.configFile = &opts.configFile
	cli.flags = opts.flags

	if cli.Config.Debug {
		debug.Enable()
	}
	if opts.Daemon {
		conf := dsshd.DSSHdConfig{}
		if err := dsshd.Start(&conf); err != nil {
			logrus.Error(err)
		}

	} else {

	}
	logrus.Info("Starting up")

	return nil
}

// configureDaemonLogs sets the logrus logging level and formatting
func configureDaemonLogs(conf *config.Config) error {
	if conf.LogLevel != "" {
		lvl, err := logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			return fmt.Errorf("unable to parse logging level: %s", conf.LogLevel)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		DisableColors:   conf.RawLogs,
		FullTimestamp:   true,
	})
	return nil
}
