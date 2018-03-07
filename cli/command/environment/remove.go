package environment

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

// NewRemoveCommand returns the `environment remove` subcommand
func NewRemoveCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove env1 env2 ...",
		Aliases: []string{"rm"},
		Short:   "Remove one or more environments",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := dockerCli.ConfigFile()
			for _, e := range args {
				if _, ok := conf.Environments[e]; !ok {
					return fmt.Errorf("no such environment: %s", e)
				}
				delete(conf.Environments, e)
				if conf.CurrentEnvironment == e {
					conf.CurrentEnvironment = ""
				}
			}
			return conf.Save()
		},
	}
	return cmd
}
