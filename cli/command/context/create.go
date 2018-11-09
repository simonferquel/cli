package context

import (
	"github.com/docker/cli/cli/context/common"
	"io/ioutil"
	"os"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/cli/cli/context/kubernetes"
	"github.com/docker/cli/cli/context/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createOptions struct {
	name                     string
	description              string // --description string: set the description
	defaultStackOrchestrator string
	docker                   dockerEndpointOptions
	kubernetes               kubernetedEndpointOptions
}

func (o *createOptions) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.description, "description", "", "set the description of the context")
	flags.StringVar(&o.defaultStackOrchestrator, "default-stack-orchestrator", "", "set the default orchestrator for stack operations if different to the default one, to use with this context (swarm|kubernetes|all)")
	o.docker.addFlags(flags)
	o.kubernetes.addFlags(flags)
}

type dockerEndpointOptions struct {
	host          string //--moby-host string: set moby endpoint host (same format as DOCKER_HOST)
	apiVersion    string
	ca            string //--moby-ca string: path to a CA file
	cert          string //--moby-cert string: path to a client cert file
	key           string //--moby-key string: path to a client key file
	skipTLSVerify bool
}

func (o *dockerEndpointOptions) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.host, "docker-host", "", "required: specify the docker endpoint on wich to connect")
	flags.StringVar(&o.apiVersion, "docker-api-version", "", "override negociated api version")
	flags.StringVar(&o.ca, "docker-tls-ca", "", "path to the ca file to validate docker endpoint")
	flags.StringVar(&o.cert, "docker-tls-cert", "", "path to the cert file to authenticate to the docker endpoint")
	flags.StringVar(&o.key, "docker-tls-key", "", "path to the key file to authenticate to the docker endpoint")
	flags.BoolVar(&o.skipTLSVerify, "docker-tls-skip-verify", false, "skip tls verify when connecting to the docker endpoint")
}

func (o *dockerEndpointOptions) toEndpoint(contextName string) (docker.Endpoint, error) {
	tlsData, err := common.TLSDataFromFiles(o.ca, o.cert, o.key)
	if err != nil {
		return docker.Endpoint{}, err
	}
	return docker.Endpoint{
		EndpointMeta: docker.EndpointMeta{
			EndpointMeta: common.EndpointMeta{
				ContextName:   contextName,
				Host:          o.host,
				SkipTLSVerify: o.skipTLSVerify,
			},
			APIVersion: o.apiVersion,
		},
		TLSData: tlsData,
	}, nil
}

type kubernetedEndpointOptions struct {
	server            string //--kubernetes-server string: kubernetes api server address
	ca                string //--kubernetes-ca string: path to a CA file
	cert              string //--kubernetes-cert string: path to a client cert file
	key               string //--kubernetes-key string: path to a client key file
	skipTLSVerify     bool
	defaultNamespace  string
	kubeconfigFile    string //--kubernetes-kubeconfig-file string: path to a kubernetes cli config file
	kubeconfigContext string //--kubernetes-kubeconfig-context string: name of the kubernetes cli config file context to use
}

func (o *kubernetedEndpointOptions) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.server, "kubernetes-host", "", "required: specify the docker endpoint on wich to connect")
	flags.StringVar(&o.ca, "kubernetes-tls-ca", "", "path to the ca file to validate kubernetes endpoint")
	flags.StringVar(&o.cert, "kubernetes-tls-cert", "", "path to the cert file to authenticate to the kubernetes endpoint")
	flags.StringVar(&o.key, "kubernetes-tls-key", "", "path to the key file to authenticate to the kubernetes endpoint")
	flags.BoolVar(&o.skipTLSVerify, "kubernetes-tls-skip-verify", false, "skip tls verify when connecting to the kubernetes endpoint")
	flags.StringVar(&o.defaultNamespace, "kubernetes-default-namespace", "default", "override default namespace when connecting to kubernetes endpoint")
	flags.StringVar(&o.kubeconfigFile, "kubernetes-kubeconfig", "", "path to an existing kubeconfig file")
	flags.StringVar(&o.kubeconfigContext, "kubernetes-kubeconfig-context", "", `context to use in the kubeconfig file referenced in "kubernetes-kubeconfig"`)
}

func (o *kubernetedEndpointOptions) toEndpoint(contextName string) (*kubernetes.Endpoint, error) {
	if o.kubeconfigFile != "" {
		ep, err := kubernetes.FromKubeConfig(contextName, o.kubeconfigFile, o.kubeconfigContext, o.defaultNamespace)
		if err != nil {
			return nil, err
		}
		return &ep, nil
	}
	if o.server != "" {
		tlsData, err := common.TLSDataFromFiles(o.ca, o.cert, o.key)
		if err != nil {
			return nil, err
		}
		return &kubernetes.Endpoint{
			EndpointMeta: kubernetes.EndpointMeta{
				EndpointMeta: common.EndpointMeta{
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

func (o *createOptions) process(s store.Store) error {
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
	dockerEP, err := o.docker.toEndpoint(o.name)
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
			return opts.process(dockerCli.ContextStore())
		},
	}

	opts.addFlags(cmd.Flags())
	return cmd
}
