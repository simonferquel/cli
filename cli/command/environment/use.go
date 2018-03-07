package environment

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

//NewUseCommand returns the `environment use` subcommand
func NewUseCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <environment name>",
		Short: "Select an environment to work with",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			conf := dockerCli.ConfigFile()
			if _, ok := conf.Environments[envName]; !ok {
				return fmt.Errorf("no such environment: %s", envName)
			}
			conf.CurrentEnvironment = envName
			return conf.Save()
		},
	}
	return cmd
}
