package kubernetes

import (
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/contextstore"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	kubernetesEndpoint   = "kubernetes"
	serverKey            = "server"
	skipTLSVerifyKey     = "skipTLSVerify"
	defaultNamespaceKey  = "defaultNamespace"
	kubeconfigFileKey    = "kubeconfigFile"
	kubeconfigContextKey = "kubeconfigContext"
	caFile               = "ca.pem"
	certFile             = "cert.pem"
	keyFile              = "key.pem"
)

// Context is a typed wrapper around a context-store context
type Context struct {
	Name              string
	Server            string
	SkipTLSVerify     bool
	DefaultNamespace  string
	KubeconfigFile    string
	KubeconfigContext string
}

// LoadKubernetesConfig loads the kubernetes client config from the context
func (c *Context) LoadKubernetesConfig(s contextstore.Store) (clientcmd.ClientConfig, error) {
	if c.KubeconfigFile != "" {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.KubeconfigFile},
			&clientcmd.ConfigOverrides{CurrentContext: c.KubeconfigContext, Context: clientcmdapi.Context{Namespace: c.DefaultNamespace}}), nil
	}
	cfg := clientcmdapi.NewConfig()
	cluster := clientcmdapi.NewCluster()
	cluster.Server = c.Server
	cluster.InsecureSkipTLSVerify = c.SkipTLSVerify
	authInfo := clientcmdapi.NewAuthInfo()
	tlsFiles, err := s.ListContextTLSFiles(c.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve context tls files")
	}
	if epTLSFiles, ok := tlsFiles[kubernetesEndpoint]; ok {
		for _, file := range epTLSFiles {
			data, err := s.GetContextTLSData(c.Name, kubernetesEndpoint, file)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load tls data")
			}
			switch file {
			case caFile:
				cluster.CertificateAuthorityData = data
			case certFile:
				authInfo.ClientCertificateData = data
			case keyFile:
				authInfo.ClientKeyData = data
			}
		}
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

func getMetaString(meta map[string]interface{}, key string) (string, bool) {
	v, ok := meta[key]
	if !ok {
		return "", false
	}
	r, ok := v.(string)
	return r, ok
}

func getMetaBool(meta map[string]interface{}, key string) (bool, bool) {
	v, ok := meta[key]
	if !ok {
		return false, false
	}
	r, ok := v.(bool)
	return r, ok
}

// ParseContext extract kubernetes endpoint info from current context
func ParseContext(name string, metadata contextstore.ContextMetadata) *Context {
	ep, ok := metadata.Endpoints[kubernetesEndpoint]
	if !ok {
		return nil
	}
	server, _ := getMetaString(ep, serverKey)
	skipTLSVerify, _ := getMetaBool(ep, skipTLSVerifyKey)
	kubeconfigFile, _ := getMetaString(ep, kubeconfigFileKey)
	kubeconfigContext, _ := getMetaString(ep, kubeconfigContextKey)
	defaultNamespace, _ := getMetaString(ep, defaultNamespaceKey)
	return &Context{
		Name:              name,
		Server:            server,
		SkipTLSVerify:     skipTLSVerify,
		KubeconfigFile:    kubeconfigFile,
		KubeconfigContext: kubeconfigContext,
		DefaultNamespace:  defaultNamespace,
	}
}

// SetKubenetesContextRaw set a context kubernetes endpoint
func SetKubenetesContextRaw(s contextstore.Store, name, server, defaultNamespace string, ca, cert, key []byte, skipTLSVerify bool) error {
	ctxMeta, err := s.GetContextMetadata(name)
	switch {
	case os.IsNotExist(err):
		ctxMeta = contextstore.ContextMetadata{
			Endpoints: make(map[string]contextstore.EndpointMetadata),
			Metadata:  make(map[string]interface{}),
		}
	case err != nil:
		return err
	}
	epMeta := make(contextstore.EndpointMetadata)
	epMeta[serverKey] = server
	epMeta[defaultNamespaceKey] = defaultNamespace
	epMeta[skipTLSVerifyKey] = skipTLSVerify
	ctxMeta.Endpoints[kubernetesEndpoint] = epMeta
	err = s.CreateOrUpdateContext(name, ctxMeta)
	if err != nil {
		return err
	}
	return s.ResetContextEndpointTLSMaterial(name, kubernetesEndpoint, createEnpointTLSData(ca, cert, key))
}

// SetKubenetesContextKubeconfig set the kubernetes endpoint of a namespace from an external kubeconfig file
func SetKubenetesContextKubeconfig(s contextstore.Store, name, kubeconfig, kubeContext, namespaceOverride string, embed bool) error {
	if embed {
		cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{CurrentContext: kubeContext, Context: clientcmdapi.Context{Namespace: namespaceOverride}})
		ns, _, err := cfg.Namespace()
		if err != nil {
			return err
		}
		clientcfg, err := cfg.ClientConfig()
		if err != nil {
			return err
		}
		var ca, key, cert []byte
		if clientcfg.CAFile != "" {
			ca, err = ioutil.ReadFile(clientcfg.CAFile)
			if err != nil {
				return err
			}
		} else {
			ca = clientcfg.CAData
		}
		if clientcfg.KeyFile != "" {
			key, err = ioutil.ReadFile(clientcfg.KeyFile)
			if err != nil {
				return err
			}
		} else {
			key = clientcfg.KeyData
		}
		if clientcfg.CertFile != "" {
			cert, err = ioutil.ReadFile(clientcfg.CertFile)
			if err != nil {
				return err
			}
		} else {
			cert = clientcfg.CertData
		}

		return SetKubenetesContextRaw(s, name, clientcfg.Host, ns, ca, cert, key, clientcfg.Insecure)
	}
	ctxMeta, err := s.GetContextMetadata(name)
	switch {
	case os.IsNotExist(err):
		ctxMeta = contextstore.ContextMetadata{
			Endpoints: make(map[string]contextstore.EndpointMetadata),
			Metadata:  make(map[string]interface{}),
		}
	case err != nil:
		return err
	}
	epMeta := make(contextstore.EndpointMetadata)
	epMeta[kubeconfigFileKey] = kubeconfig
	epMeta[kubeconfigContextKey] = kubeContext
	epMeta[defaultNamespaceKey] = namespaceOverride
	ctxMeta.Endpoints[kubernetesEndpoint] = epMeta
	return s.CreateOrUpdateContext(name, ctxMeta)
}

func createEnpointTLSData(ca, cert, key []byte) *contextstore.EndpointTLSData {
	if ca == nil && cert == nil && key == nil {
		return nil
	}
	result := contextstore.EndpointTLSData{
		Files: make(map[string][]byte),
	}
	if ca != nil {
		result.Files[caFile] = ca
	}
	if cert != nil {
		result.Files[certFile] = cert
	}
	if key != nil {
		result.Files[keyFile] = key
	}
	return &result
}
