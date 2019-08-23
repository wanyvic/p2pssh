package main

import (
	"github.com/docker/docker/opts"
	"github.com/spf13/pflag"
	"github.com/wanyvic/dssh/cmd/config"
)

type MassgridRPCOptions struct {
	MassGridHost []string
	Username     string
	Password     string
}

type daemonOptions struct {
	configFile   string
	daemonConfig *config.Config
	flags        *pflag.FlagSet
	Debug        bool
	Hosts        []string
	LogLevel     string
	Relay        bool
	Daemon       bool
}

// newDaemonOptions returns a new daemonFlags
func newDaemonOptions(config *config.Config) *daemonOptions {
	return &daemonOptions{
		daemonConfig: config,
	}
}

// InstallFlags adds flags for the common options on the FlagSet
func (o *daemonOptions) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&o.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	flags.BoolVarP(&o.Relay, "relay", "R", true, "Enable relay mode")
	flags.BoolVarP(&o.Daemon, "daemon", "", false, "Enable daemon mode")

	hostOpt := opts.NewNamedListOptsRef("hosts", &o.Hosts, opts.ValidateHost)
	flags.VarP(hostOpt, "host", "H", "Daemon socket(s) to connect to")

}

// SetDefaultOptions sets default values for options after flag parsing is
// complete
func (o *daemonOptions) SetDefaultOptions(flags *pflag.FlagSet) {
	o.Relay = true
	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on TLS
	// TLSVerify can be true even if not set due to DOCKER_TLS_VERIFY env var, so we need
	// to check that here as well

}
