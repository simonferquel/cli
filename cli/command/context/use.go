package context

import (
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func newUseCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use <context name>",
		Aliases: []string{"select", "switch"},
		Short:   "Set the current docker context",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if _, err := dockerCli.ContextStore().GetContextMetadata(name); err != nil {
				return err
			}
			return dockerCli.ContextStore().SetCurrentContext(name)
		},
	}
	return cmd
}
