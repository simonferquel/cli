package context

import (
	"io/ioutil"
	"os"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/kubernetes"
	"github.com/docker/context-store"
	dockerCtx "github.com/docker/docker/client/context"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createOptions struct {
	name                        string
	description                 string // --description string: set the description
	defaultOrchestrator         string //--default-orchestrator <kubernetes|swarm|all>: set default orchestrator for orchestrator-specific commands
	defaultStackOrchestrator    string
	dockerHost                  string //--moby-host string: set moby endpoint host (same format as DOCKER_HOST)
	dockerAPIVersion            string
	dockerCA                    string //--moby-ca string: path to a CA file
	dockerCert                  string //--moby-cert string: path to a client cert file
	dockerKey                   string //--moby-key string: path to a client key file
	dockerSkipTLSVerify         bool
	kubernetesServer            string //--kubernetes-server string: kubernetes api server address
	kubernetesCA                string //--kubernetes-ca string: path to a CA file
	kubernetesCert              string //--kubernetes-cert string: path to a client cert file
	kubernetesKey               string //--kubernetes-key string: path to a client key file
	kubernetesSkipTLSVerify     bool
	kubernetesDefaultNamespace  string
	kubernetesKubeconfigFile    string //--kubernetes-kubeconfig-file string: path to a kubernetes cli config file
	kubernetesKubeconfigContext string //--kubernetes-kubeconfig-context string: name of the kubernetes cli config file context to use
	kubernetesKubeconfigEmbed   bool   //--kubernetes-embed bool (default: false): only used with --kubernetes-kubeconfig-file / --kubernetes-kubeconfig-context: embed Kubernetes configuration in the context store instead of referencing the config file
}

func (o *createOptions) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.description, "description", "", "set the description of the context")
	flags.StringVar(&o.defaultOrchestrator, "default-orchestrator", "swarm", "set the default orchestrator to use with this context (swarm|kubernetes|all)")
	flags.StringVar(&o.defaultStackOrchestrator, "default-stack-orchestrator", "", "set the default orchestrator for stack operations if different to the default one, to use with this context (swarm|kubernetes|all)")
	flags.StringVar(&o.dockerHost, "docker-host", "", "required: specify the docker endpoint on wich to connect")
	flags.StringVar(&o.dockerAPIVersion, "docker-api-version", "", "override negociated api version")
	flags.StringVar(&o.dockerCA, "docker-tls-ca", "", "path to the ca file to validate docker endpoint")
	flags.StringVar(&o.dockerCert, "docker-tls-cert", "", "path to the cert file to authenticate to the docker endpoint")
	flags.StringVar(&o.dockerKey, "docker-tls-key", "", "path to the key file to authenticate to the docker endpoint")
	flags.BoolVar(&o.dockerSkipTLSVerify, "docker-tls-skip-verify", false, "skip tls verify when connecting to the docker endpoint")
	flags.StringVar(&o.kubernetesServer, "kubernetes-host", "", "required: specify the docker endpoint on wich to connect")
	flags.StringVar(&o.kubernetesCA, "kubernetes-tls-ca", "", "path to the ca file to validate kubernetes endpoint")
	flags.StringVar(&o.kubernetesCert, "kubernetes-tls-cert", "", "path to the cert file to authenticate to the kubernetes endpoint")
	flags.StringVar(&o.kubernetesKey, "kubernetes-tls-key", "", "path to the key file to authenticate to the kubernetes endpoint")
	flags.BoolVar(&o.kubernetesSkipTLSVerify, "kubernetes-tls-skip-verify", false, "skip tls verify when connecting to the kubernetes endpoint")
	flags.StringVar(&o.kubernetesDefaultNamespace, "kubernetes-default-namespace", "default", "override default namespace when connecting to kubernetes endpoint")
	flags.StringVar(&o.kubernetesKubeconfigFile, "kubernetes-kubeconfig", "", "path to an existing kubeconfig file")
	flags.StringVar(&o.kubernetesKubeconfigContext, "kubernetes-kubeconfig-context", "", `context to use in the kubeconfig file referenced in "kubernetes-kubeconfig"`)
	flags.BoolVar(&o.kubernetesKubeconfigEmbed, "kubernetes-kubeconfig-embed", false, `if kubernetes-kubeconfig is specified, embed the config in the docker context store instead of referencing it`)
}

func loadFileIfNotEmpty(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	return ioutil.ReadFile(path)
}

func (o *createOptions) process(s store.Store) error {
	if _, err := s.GetContextMetadata(o.name); !os.IsNotExist(err) {
		if err != nil {
			return errors.Wrap(err, "error while getting existing contexts")
		}
		return errors.Errorf("context %q already exists", o.name)
	}
	defaultOrchestrator, err := command.NormalizeOrchestrator(o.defaultOrchestrator)
	if err != nil {
		return errors.Wrap(err, "unable to parse default-orchestrator")
	}
	stackOrchestrator, err := command.NormalizeOrchestrator(o.defaultStackOrchestrator)
	if err != nil {
		return errors.Wrap(err, "unable to parse default-stack-orchestrator")
	}
	dockerCA, err := loadFileIfNotEmpty(o.dockerCA)
	if err != nil {
		return errors.Wrap(err, "unable to load docker-tls-ca")
	}
	dockerCert, err := loadFileIfNotEmpty(o.dockerCert)
	if err != nil {
		return errors.Wrap(err, "unable to load docker-tls-cert")
	}
	dockerKey, err := loadFileIfNotEmpty(o.dockerKey)
	if err != nil {
		return errors.Wrap(err, "unable to load docker-tls-key")
	}
	if err = dockerCtx.SetDockerEndpoint(s, o.name, o.dockerHost, o.dockerAPIVersion, dockerCA, dockerCert, dockerKey, o.dockerSkipTLSVerify); err != nil {
		return errors.Wrap(err, "unable to set docker endpoint")
	}

	if o.kubernetesKubeconfigFile != "" {
		if err = kubernetes.SetKubenetesContextKubeconfig(s, o.name, o.kubernetesKubeconfigFile, o.kubernetesKubeconfigContext, o.kubernetesKubeconfigEmbed); err != nil {
			return errors.Wrap(err, "unable to set kubernetes endpoint")
		}
	} else if o.kubernetesServer != "" {
		kubeCA, err := loadFileIfNotEmpty(o.kubernetesCA)
		if err != nil {
			return errors.Wrap(err, "unable to load kubernetes-tls-ca")
		}
		kubeCert, err := loadFileIfNotEmpty(o.kubernetesCert)
		if err != nil {
			return errors.Wrap(err, "unable to load kubernetes-tls-cert")
		}
		kubeKey, err := loadFileIfNotEmpty(o.kubernetesKey)
		if err != nil {
			return errors.Wrap(err, "unable to load kubernetes-tls-key")
		}
		if err = kubernetes.SetKubenetesContextRaw(s, o.name, o.kubernetesServer, o.kubernetesDefaultNamespace,
			kubeCA, kubeCert, kubeKey, o.kubernetesSkipTLSVerify); err != nil {
			return errors.Wrap(err, "unable to set kubernetes endpoint")
		}
	}
	// at this point, the context should exist with endpoint configurations
	ctx, err := s.GetContextMetadata(o.name)
	if err != nil {
		return errors.Wrap(err, "error while getting context")
	}
	command.SetContextMetadata(&ctx, command.ContextMetadata{
		Description:       o.description,
		Orchestrator:      defaultOrchestrator,
		StackOrchestrator: stackOrchestrator,
	})

	return s.CreateOrUpdateContext(o.name, ctx)
}

func newCreateCommand(dockerCli command.Cli) *cobra.Command {
	opts := &createOptions{}
	cmd := &cobra.Command{
		Use:   "create <name> [options]",
		Short: "create a context",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			return opts.process(dockerCli.ContextStore())
		},
	}

	opts.addFlags(cmd.Flags())
	return cmd
}
