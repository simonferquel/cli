package context

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/cli/cli/context/kubernetes"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/docker/pkg/homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createOptions struct {
	name                     string
	description              string // --description string: set the description
	defaultStackOrchestrator string
	docker                   dockerEndpointOptions
	kubernetes               kubernetesEndpointOptions
}

func (o *createOptions) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.description, "description", "", "set the description of the context")
	flags.StringVar(&o.defaultStackOrchestrator, "default-stack-orchestrator", "", "set the default orchestrator for stack operations if different to the default one, to use with this context (swarm|kubernetes|all)")
	o.docker.addFlags(flags, "docker-")
	o.kubernetes.addFlags(flags, "kubernetes-")
}

type dockerEndpointOptions struct {
	host          string //--moby-host string: set moby endpoint host (same format as DOCKER_HOST)
	apiVersion    string
	ca            string //--moby-ca string: path to a CA file
	cert          string //--moby-cert string: path to a client cert file
	key           string //--moby-key string: path to a client key file
	skipTLSVerify bool
	fromEnv       bool
}

func (o *dockerEndpointOptions) addFlags(flags *pflag.FlagSet, prefix string) {
	flags.StringVar(&o.host, prefix+"host", "", "required: specify the docker endpoint on wich to connect")
	flags.StringVar(&o.apiVersion, prefix+"api-version", "", "override negociated api version")
	flags.StringVar(&o.ca, prefix+"tls-ca", "", "path to the ca file to validate docker endpoint")
	flags.StringVar(&o.cert, prefix+"tls-cert", "", "path to the cert file to authenticate to the docker endpoint")
	flags.StringVar(&o.key, prefix+"tls-key", "", "path to the key file to authenticate to the docker endpoint")
	flags.BoolVar(&o.skipTLSVerify, prefix+"tls-skip-verify", false, "skip tls verify when connecting to the docker endpoint")
	flags.BoolVar(&o.fromEnv, prefix+"from-env", false, "convert the current env-variable based configuration to a context")
}

func (o *dockerEndpointOptions) toEndpoint(cli command.Cli, contextName string) (docker.Endpoint, error) {
	if o.fromEnv {
		if cli.CurrentContext() != command.ContextDockerHost {
			return docker.Endpoint{}, errors.New("cannot create a context from environment when a context is in use")
		}
		ep := cli.DockerEndpoint()
		ep.ContextName = contextName
		return ep, nil
	}
	tlsData, err := context.TLSDataFromFiles(o.ca, o.cert, o.key)
	if err != nil {
		return docker.Endpoint{}, err
	}
	return docker.Endpoint{
		EndpointMeta: docker.EndpointMeta{
			EndpointMetaBase: context.EndpointMetaBase{
				ContextName:   contextName,
				Host:          o.host,
				SkipTLSVerify: o.skipTLSVerify,
			},
			APIVersion: o.apiVersion,
		},
		TLSData: tlsData,
	}, nil
}

type kubernetesEndpointOptions struct {
	server            string //--kubernetes-server string: kubernetes api server address
	ca                string //--kubernetes-ca string: path to a CA file
	cert              string //--kubernetes-cert string: path to a client cert file
	key               string //--kubernetes-key string: path to a client key file
	skipTLSVerify     bool
	defaultNamespace  string
	kubeconfigFile    string //--kubernetes-kubeconfig-file string: path to a kubernetes cli config file
	kubeconfigContext string //--kubernetes-kubeconfig-context string: name of the kubernetes cli config file context to use
	fromEnv           bool
}

func (o *kubernetesEndpointOptions) addFlags(flags *pflag.FlagSet, prefix string) {
	flags.StringVar(&o.server, prefix+"host", "", "specify the kubernetes endpoint on wich to connect")
	flags.StringVar(&o.ca, prefix+"tls-ca", "", "path to the ca file to validate kubernetes endpoint")
	flags.StringVar(&o.cert, prefix+"tls-cert", "", "path to the cert file to authenticate to the kubernetes endpoint")
	flags.StringVar(&o.key, prefix+"tls-key", "", "path to the key file to authenticate to the kubernetes endpoint")
	flags.BoolVar(&o.skipTLSVerify, prefix+"tls-skip-verify", false, "skip tls verify when connecting to the kubernetes endpoint")
	flags.StringVar(&o.defaultNamespace, prefix+"default-namespace", "default", "override default namespace when connecting to kubernetes endpoint")
	flags.StringVar(&o.kubeconfigFile, prefix+"kubeconfig", "", "path to an existing kubeconfig file")
	flags.StringVar(&o.kubeconfigContext, prefix+"kubeconfig-context", "", `context to use in the kubeconfig file referenced in "kubernetes-kubeconfig"`)
	flags.BoolVar(&o.fromEnv, prefix+"from-env", false, `use the default kubeconfig file or the value defined in KUBECONFIG environement variable`)
}

func (o *kubernetesEndpointOptions) toEndpoint(contextName string) (*kubernetes.Endpoint, error) {
	if o.kubeconfigFile == "" && o.fromEnv {
		if config := os.Getenv("KUBECONFIG"); config != "" {
			o.kubeconfigFile = config
		} else {
			o.kubeconfigFile = filepath.Join(homedir.Get(), ".kube/config")
		}
	}
	if o.kubeconfigFile != "" {
		ep, err := kubernetes.FromKubeConfig(contextName, o.kubeconfigFile, o.kubeconfigContext, o.defaultNamespace)
		if err != nil {
			return nil, err
		}
		return &ep, nil
	}
	if o.server != "" {
		tlsData, err := context.TLSDataFromFiles(o.ca, o.cert, o.key)
		if err != nil {
			return nil, err
		}
		return &kubernetes.Endpoint{
			EndpointMeta: kubernetes.EndpointMeta{
				EndpointMetaBase: context.EndpointMetaBase{
					ContextName:   contextName,
					Host:          o.server,
					SkipTLSVerify: o.skipTLSVerify,
				},
				DefaultNamespace: o.defaultNamespace,
			},
			TLSData: tlsData,
		}, nil
	}
	return nil, nil
}

func loadFileIfNotEmpty(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	return ioutil.ReadFile(path)
}

func (o *createOptions) process(cli command.Cli, s store.Store) error {
	if _, err := s.GetContextMetadata(o.name); !os.IsNotExist(err) {
		if err != nil {
			return errors.Wrap(err, "error while getting existing contexts")
		}
		return errors.Errorf("context %q already exists", o.name)
	}
	stackOrchestrator, err := command.NormalizeOrchestrator(o.defaultStackOrchestrator)
	if err != nil {
		return errors.Wrap(err, "unable to parse default-stack-orchestrator")
	}
	dockerEP, err := o.docker.toEndpoint(cli, o.name)
	if err != nil {
		return errors.Wrap(err, "unable to create docker endpoint config")
	}
	if err := docker.Save(s, dockerEP); err != nil {
		return errors.Wrap(err, "unable to save docker endpoint config")
	}
	kubernetesEP, err := o.kubernetes.toEndpoint(o.name)
	if err != nil {
		return errors.Wrap(err, "unable to create kubernetes endpoint config")
	}
	if kubernetesEP != nil {
		if err := kubernetes.Save(s, *kubernetesEP); err != nil {
			return errors.Wrap(err, "unable to save kubernetes endpoint config")
		}
	}

	// at this point, the context should exist with endpoints configuration
	ctx, err := s.GetContextMetadata(o.name)
	if err != nil {
		return errors.Wrap(err, "error while getting context")
	}
	command.SetContextMetadata(&ctx, command.ContextMetadata{
		Description:       o.description,
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
			return opts.process(dockerCli, dockerCli.ContextStore())
		},
	}

	opts.addFlags(cmd.Flags())
	return cmd
}
