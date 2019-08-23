package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/dssh/cmd/cli"
	"github.com/wanyvic/dssh/cmd/config"
	"github.com/wanyvic/dssh/dsshversion"
)

func newDaemonCommand() (*cobra.Command, error) {
	opts := newDaemonOptions(config.New())

	cmd := &cobra.Command{
		Use:           "pssh [OPTIONS]",
		Short:         "A Distributed Secure Shell",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			if cmd.HasSubCommands() {
				return errors.Errorf("\n" + strings.TrimRight(cmd.UsageString(), "\n"))
			}

			return errors.Errorf(
				"\"%s\" accepts no argument(s).\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
				cmd.CommandPath(),
				cmd.CommandPath(),
				cmd.UseLine(),
				cmd.Short,
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.flags = cmd.Flags()
			return runDaemon(opts)
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", dsshversion.Version, dsshversion.GitCommit),
	}
	cli.SetupRootCommand(cmd)
	command.AddCommand(cmd)

	flags := cmd.Flags()
	flags.BoolP("version", "v", false, "Print version information and quit")
	defaultDaemonConfigFile, err := getDefaultDaemonConfigFile()
	if err != nil {
		return nil, err
	}
	flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "Daemon configuration file")
	opts.InstallFlags(flags)
	if err := installConfigFlags(opts.daemonConfig, flags); err != nil {
		return nil, err
	}
	return cmd, nil
}

func main() {

	// initial log formatting; this setting is updated after the daemon configuration is loaded.
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
		FullTimestamp:   true,
	})

	cmd, err := newDaemonCommand()
	if err != nil {
		panic(err)
	}
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
	logrus.Debug("pssh exit")
}
