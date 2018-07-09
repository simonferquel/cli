package context

import (
	"fmt"
	"sort"

	"vbom.ml/util/sortorder"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/formatter"
	"github.com/docker/cli/kubernetes"
	dcontext "github.com/docker/docker/client/context"
	"github.com/spf13/cobra"
)

type listOptions struct {
	format string
}

func newListCommand(dockerCli command.Cli) *cobra.Command {
	opts := &listOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List contexts",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.format, "format", "", "Pretty-print contexts using a Go template")
	return cmd
}

func runList(dockerCli command.Cli, opts *listOptions) error {
	curContext := dockerCli.CurrentContext()
	contextMap, err := dockerCli.ContextStore().ListContexts()
	if err != nil {
		return err
	}
	var contexts []*formatter.ClientContext
	for name, rawMeta := range contextMap {
		meta, err := command.GetContextMetadata(rawMeta)
		if err != nil {
			return err
		}
		dockerEndpoint, err := dcontext.Parse(name, rawMeta)
		if err != nil {
			return err
		}
		kubernetesEndpoint := kubernetes.ParseContext(name, rawMeta)
		kubEndpointText := ""
		if kubernetesEndpoint != nil {
			if kubernetesEndpoint.KubeconfigFile != "" {
				kubEndpointText = fmt.Sprintf("%s (%s)", kubernetesEndpoint.KubeconfigFile, kubernetesEndpoint.KubeconfigContext)
			} else {
				kubEndpointText = fmt.Sprintf("%s (%s)", kubernetesEndpoint.Server, kubernetesEndpoint.DefaultNamespace)
			}
		}
		desc := formatter.ClientContext{
			Name:               name,
			Current:            name == curContext,
			Description:        meta.Description,
			Orchestrator:       string(meta.Orchestrator),
			StackOrchestrator:  string(meta.StackOrchestrator),
			DockerEndpoint:     dockerEndpoint.Host,
			KubernetesEndpoint: kubEndpointText,
		}
		contexts = append(contexts, &desc)
	}
	sort.Slice(contexts, func(i, j int) bool {
		return sortorder.NaturalLess(contexts[i].Name, contexts[j].Name)
	})
	return format(dockerCli, opts, contexts)
}

func format(dockerCli command.Cli, opts *listOptions, contexts []*formatter.ClientContext) error {
	format := opts.format
	if format == "" {
		format = formatter.ClientContextTableFormat
	}
	contextCtx := formatter.Context{
		Output: dockerCli.Out(),
		Format: formatter.Format(format),
	}
	return formatter.ClientContextWrite(contextCtx, contexts)
}
