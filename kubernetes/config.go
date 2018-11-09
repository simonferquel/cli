package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/kubernetes"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/docker/pkg/homedir"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubernetesConfig resolves the path to the desired Kubernetes configuration file based on
// the KUBECONFIG environment variable and command line flags.
func NewKubernetesConfig(s store.Store, contextName, configPath string) (clientcmd.ClientConfig, error) {
	if configPath == "" && contextName != command.ContextDockerHost {
		ctxMeta, err := s.GetContextMetadata(contextName)
		if err != nil {
			return nil, err
		}
		epMeta := kubernetes.Parse(contextName, ctxMeta)
		if epMeta != nil {
			ep, err := epMeta.WithTLSData(s)
			if err != nil {
				return nil, err
			}
			return ep.KubernetesConfig()
		}
	}
	kubeConfig := configPath
	if kubeConfig == "" {
		if config := os.Getenv("KUBECONFIG"); config != "" {
			kubeConfig = config
		} else {
			kubeConfig = filepath.Join(homedir.Get(), ".kube/config")
		}
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{}), nil
}
