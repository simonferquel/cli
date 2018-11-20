package kubernetes

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/cli/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// EndpointMeta is a typed wrapper around a context-store generic endpoint describing
// a Kubernetes endpoint, without TLS data
type EndpointMeta struct {
	context.EndpointMetaBase
	DefaultNamespace string
}

// Endpoint is a typed wrapper around a context-store generic endpoint describing
// a Kubernetes endpoint, with TLS data
type Endpoint struct {
	EndpointMeta
	TLSData *context.TLSData
}

// WithTLSData loads TLS materials for the endpoint
func (c *EndpointMeta) WithTLSData(s store.Store) (Endpoint, error) {
	tlsData, err := context.LoadTLSData(s, c.ContextName, kubernetesEndpointKey)
	if err != nil {
		return Endpoint{}, err
	}
	return Endpoint{
		EndpointMeta: *c,
		TLSData:      tlsData,
	}, nil
}

// KubernetesConfig creates the kubernetes client config from the endpoint
func (c *Endpoint) KubernetesConfig() (clientcmd.ClientConfig, error) {
	cfg := clientcmdapi.NewConfig()
	cluster := clientcmdapi.NewCluster()
	cluster.Server = c.Host
	cluster.InsecureSkipTLSVerify = c.SkipTLSVerify
	authInfo := clientcmdapi.NewAuthInfo()
	if c.TLSData != nil {
		cluster.CertificateAuthorityData = c.TLSData.CA
		authInfo.ClientCertificateData = c.TLSData.Cert
		authInfo.ClientKeyData = c.TLSData.Key
	}
	cfg.Clusters["cluster"] = cluster
	cfg.AuthInfos["authInfo"] = authInfo
	ctx := clientcmdapi.NewContext()
	ctx.AuthInfo = "authInfo"
	ctx.Cluster = "cluster"
	ctx.Namespace = c.DefaultNamespace
	cfg.Contexts["context"] = ctx
	cfg.CurrentContext = "context"
	return clientcmd.NewDefaultClientConfig(*cfg, &clientcmd.ConfigOverrides{}), nil
}

// EndpointFromContext extracts kubernetes endpoint info from current context
func EndpointFromContext(name string, metadata store.ContextMetadata) *EndpointMeta {
	ep, ok := metadata.Endpoints[kubernetesEndpointKey]
	if !ok {
		return nil
	}
	commonMeta := context.EndpointFromContext(name, kubernetesEndpointKey, metadata)
	if commonMeta == nil {
		return nil
	}
	defaultNamespace, _ := ep.GetString(defaultNamespaceKey)
	return &EndpointMeta{
		EndpointMetaBase: *commonMeta,
		DefaultNamespace: defaultNamespace,
	}
}

// ConfigFromContext resolves a kubernetes client config for the specified context.
// If kubeconfigOverride is specified, use this config file instead of the context defaults.ConfigFromContext
// if command.ContextDockerHost is specified as the context name, fallsback to the default user's kubeconfig file
func ConfigFromContext(name string, s store.Store, kubeconfigOverride string) (clientcmd.ClientConfig, error) {
	if kubeconfigOverride != "" || name == command.ContextDockerHost {
		return kubernetes.NewKubernetesConfig(kubeconfigOverride)
	}

	ctxMeta, err := s.GetContextMetadata(name)
	if err != nil {
		return nil, err
	}
	epMeta := EndpointFromContext(name, ctxMeta)
	if epMeta != nil {
		ep, err := epMeta.WithTLSData(s)
		if err != nil {
			return nil, err
		}
		return ep.KubernetesConfig()
	}
	// context has no kubernetes endpoint
	return kubernetes.NewKubernetesConfig("")
}
