package context

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/kubernetes"
	"github.com/docker/cli/cli/context/store"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type exportOptions struct {
	kubeconfig  bool
	contextName string
	dest        string
}

func newExportCommand(dockerCli command.Cli) *cobra.Command {
	opts := &exportOptions{}
	cmd := &cobra.Command{
		Use:   "export <context> [output file]",
		Short: "Export a context",
		Args:  cli.RequiresRangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.contextName = args[0]
			if len(args) == 2 {
				opts.dest = args[1]
			} else {
				opts.dest = opts.contextName
				if opts.kubeconfig {
					opts.dest += ".kubeconfig"
				} else {
					opts.dest += ".dockercontext"
				}
			}
			return runExport(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.kubeconfig, "kubeconfig", "k", false, "export as a kubeconfig file")
	return cmd
}
func runExport(dockerCli command.Cli, opts *exportOptions) error {
	ctxMeta, err := dockerCli.ContextStore().GetContextMetadata(opts.contextName)
	if err != nil {
		return err
	}
	if !opts.kubeconfig {
		reader := store.Export(opts.contextName, dockerCli.ContextStore())
		defer reader.Close()
		f, err := os.OpenFile(opts.dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, reader)
		return err
	}
	kubernetesEndpointMeta := kubernetes.Parse(opts.contextName, ctxMeta)
	if kubernetesEndpointMeta == nil {
		return fmt.Errorf("context %q has no kubernetes endpoint", opts.contextName)
	}
	kubernetesEndpoint, err := kubernetesEndpointMeta.WithTLSData(dockerCli.ContextStore())
	if err != nil {
		return err
	}
	kubeConfig, err := kubernetesEndpoint.KubernetesConfig()
	if err != nil {
		return err
	}
	rawCfg, err := kubeConfig.RawConfig()
	if err != nil {
		return err
	}
	return clientcmd.WriteToFile(rawCfg, opts.dest)
}
