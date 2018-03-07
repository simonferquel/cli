package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/spf13/cobra"
)

type importOpts struct {
	name      string
	overwrite bool
	use       bool
}

//NewImportCommand returns the `environment import` subcommand
func NewImportCommand(dockerCli command.Cli) *cobra.Command {
	opts := &importOpts{}
	cmd := &cobra.Command{
		Use:   "import --name <name> <path.dockerenv>",
		Short: "Import an environment",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.name == "" {
				return errors.New("name is required")
			}
			conf := dockerCli.ConfigFile()
			if _, ok := conf.Environments[opts.name]; ok && !opts.overwrite {
				return fmt.Errorf("environment %s already exists", opts.name)
			}
			bytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			env := configfile.Environment{}
			if err := json.Unmarshal(bytes, &env); err != nil {
				return err
			}
			conf.Environments[opts.name] = env
			if opts.use {
				conf.CurrentEnvironment = opts.name
			}
			return conf.Save()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.name, "name", "n", "", "environment name")
	flags.BoolVar(&opts.overwrite, "force", false, "overwrite existing environment")
	flags.BoolVar(&opts.use, "use", false, "use this environment as your current environment")
	return cmd
}
