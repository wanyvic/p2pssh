package cli

import (
	"github.com/docker/docker-ce/components/cli/cli/command"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

// NewContainerCommand returns a cobra command for `container` subcommands
func NewContainerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Manage containers",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(dockerCli.Err()),
	}
	cmd.AddCommand()
	return cmd
}
