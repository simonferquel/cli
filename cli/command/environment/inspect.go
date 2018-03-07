package environment

import (
	"fmt"

	"github.com/docker/cli/cli/command/inspect"
	"github.com/docker/cli/cli/config/configfile"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	format string
}

type namedEnv struct {
	Name                   string `json:"name,omitempty"`
	configfile.Environment `json:","`
}

// NewInspectCommand returns the `environment inspect` subcommand
func NewInspectCommand(dockerCli command.Cli) *cobra.Command {
	opts := &inspectOptions{}
	cmd := &cobra.Command{
		Use:   "inspect env1...",
		Short: "Inspect environments",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := dockerCli.ConfigFile()
			return inspect.Inspect(dockerCli.Out(), args, opts.format, func(ref string) (interface{}, []byte, error) {
				env, ok := conf.Environments[ref]
				if !ok {
					return nil, nil, fmt.Errorf("no such environment: %s", ref)
				}
				return namedEnv{Name: ref, Environment: env}, nil, nil
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")
	return cmd
}
