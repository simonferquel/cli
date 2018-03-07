package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

//NewExportCommand returns the `environment export` subcommand
func NewExportCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name> [dir/[filename.dockerenv]]",
		Short: "Export an environment",
		Args:  cli.RequiresRangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0] + ".dockerenv"
			if len(args) == 2 {
				target = args[1]
				s, err := os.Stat(target)
				if err == nil && s.IsDir() {
					target = filepath.Join(target, args[0]+".dockerenv")
				} else if err != nil && !os.IsNotExist(err) {
					return err
				}
			}
			env, ok := dockerCli.ConfigFile().Environments[args[0]]
			if !ok {
				return fmt.Errorf("no such environment: %s", args[0])
			}
			bytes, err := json.Marshal(&env)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(target, bytes, 0644)
		},
	}
	return cmd
}
