package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/docker/pkg/homedir"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// NewKubernetesConfig resolves the path to the desired Kubernetes configuration file, depending
// environment variable and command line flag.
func NewKubernetesConfig(configFlag string, env *configfile.Environment) (*restclient.Config, error) {
	if env != nil && env.Kubernetes != nil {
		if env.Kubernetes.KubeconfigContext != "" || env.Kubernetes.KubeconfigFile != "" {
			kubeConfig := env.Kubernetes.KubeconfigFile
			if kubeConfig == "" {
				kubeConfig = filepath.Join(homedir.Get(), ".kube/config")
			}
			return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
				&clientcmd.ConfigOverrides{CurrentContext: env.Kubernetes.KubeconfigContext},
			).ClientConfig()
		}
		cfg := api.NewConfig()
		auth := api.NewAuthInfo()
		auth.ClientCertificate = env.Kubernetes.Cert
		auth.ClientCertificateData = env.Kubernetes.CertData
		auth.ClientKey = env.Kubernetes.Key
		auth.ClientKeyData = env.Kubernetes.KeyData
		cluster := api.NewCluster()
		cluster.CertificateAuthority = env.Kubernetes.Ca
		cluster.CertificateAuthorityData = env.Kubernetes.CaData
		cluster.InsecureSkipTLSVerify = env.Kubernetes.SkipTLSVerify
		cluster.Server = env.Kubernetes.Server
		ctx := api.NewContext()
		ctx.AuthInfo = "default"
		ctx.Cluster = "default"
		cfg.AuthInfos["default"] = auth
		cfg.Clusters["default"] = cluster
		cfg.Contexts["default"] = ctx
		cfg.CurrentContext = "default"
		return clientcmd.NewDefaultClientConfig(*cfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	}
	kubeConfig := configFlag
	if kubeConfig == "" {
		if config := os.Getenv("KUBECONFIG"); config != "" {
			kubeConfig = config
		} else {
			kubeConfig = filepath.Join(homedir.Get(), ".kube/config")
		}
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfig)
}
