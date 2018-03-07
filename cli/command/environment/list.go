package environment

import (
	"fmt"
	"sort"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

// NewListCommand returns the `environment list` subcommand
func NewListCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List environments",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := dockerCli.ConfigFile()
			sortedEnvs := make([]string, len(conf.Environments))
			ix := 0
			for k := range conf.Environments {
				sortedEnvs[ix] = k
				ix++
			}
			sort.Strings(sortedEnvs)
			for _, e := range sortedEnvs {
				if e == conf.CurrentEnvironment {
					fmt.Fprintln(dockerCli.Out(), e, " *")
				} else {
					fmt.Fprintln(dockerCli.Out(), e)
				}
			}
			return nil
		},
	}
	return cmd
}
