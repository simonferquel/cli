package environment

import (
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
)

// NewEnvironmentCommand returns a cobra command for `environment` subcommands
func NewEnvironmentCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment",
		Aliases: []string{"env"},
		Short:   "Manage environments",
		Args:    cli.NoArgs,
		RunE:    command.ShowHelp(dockerCli.Err()),
	}
	cmd.AddCommand(
		NewListCommand(dockerCli),
		NewUseCommand(dockerCli),
		NewExportCommand(dockerCli),
		NewRemoveCommand(dockerCli),
		NewImportCommand(dockerCli),
		NewInspectCommand(dockerCli),
	)

	return cmd
}
