package kubernetes

import (
	"github.com/docker/cli/cli/context/common"
	"github.com/docker/cli/cli/context/store"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// EndpointMeta is a typed wrapper around a context-store generic endpoint describing
// a Kubernetes endpoint, without TLS data
type EndpointMeta struct {
	common.EndpointMeta
	DefaultNamespace string
}

// Endpoint is a typed wrapper around a context-store generic endpoint describing
// a Kubernetes endpoint, with TLS data
type Endpoint struct {
	EndpointMeta
	TLSData *common.TLSData
}

// WithTLSData loads TLS materials for the endpoint
func (c *EndpointMeta) WithTLSData(s store.Store) (Endpoint, error) {
	tlsData, err := common.LoadTLSData(s, c.ContextName, kubernetesEndpointKey)
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

// Parse extracts kubernetes endpoint info from current context
func Parse(name string, metadata store.ContextMetadata) *EndpointMeta {
	ep, ok := metadata.Endpoints[kubernetesEndpointKey]
	if !ok {
		return nil
	}
	commonMeta := common.Parse(name, kubernetesEndpointKey, metadata)
	if commonMeta == nil {
		return nil
	}
	defaultNamespace, _ := ep.GetString(defaultNamespaceKey)
	return &EndpointMeta{
		EndpointMeta:     *commonMeta,
		DefaultNamespace: defaultNamespace,
	}
}
