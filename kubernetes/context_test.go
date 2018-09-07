package kubernetes

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/pkg/contextstore"
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestSaveLoadContexts(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "test-load-save-k8-context")
	assert.Check(t, err)
	defer os.RemoveAll(storeDir)
	store, err := contextstore.NewStore(storeDir)
	assert.Check(t, err)
	assert.Check(t, SetKubenetesContextRaw(store, "raw-notls", "https://test", "test", nil, nil, nil, false))
	assert.Check(t, SetKubenetesContextRaw(store, "raw-notls-skip", "https://test", "test", nil, nil, nil, true))
	assert.Check(t, SetKubenetesContextRaw(store, "raw-tls", "https://test", "test", []byte("ca"), []byte("cert"), []byte("key"), true))

	kcFile, err := ioutil.TempFile(os.TempDir(), "test-load-save-k8-context")
	assert.Check(t, err)
	defer os.Remove(kcFile.Name())
	defer kcFile.Close()
	cfg := clientcmdapi.NewConfig()
	cfg.AuthInfos["user"] = clientcmdapi.NewAuthInfo()
	cfg.Contexts["context1"] = clientcmdapi.NewContext()
	cfg.Clusters["cluster1"] = clientcmdapi.NewCluster()
	cfg.Contexts["context2"] = clientcmdapi.NewContext()
	cfg.Clusters["cluster2"] = clientcmdapi.NewCluster()
	cfg.AuthInfos["user"].ClientCertificateData = []byte("cert")
	cfg.AuthInfos["user"].ClientKeyData = []byte("key")
	cfg.Clusters["cluster1"].Server = "https://server1"
	cfg.Clusters["cluster1"].InsecureSkipTLSVerify = true
	cfg.Clusters["cluster2"].Server = "https://server2"
	cfg.Clusters["cluster2"].CertificateAuthorityData = []byte("ca")
	cfg.Contexts["context1"].AuthInfo = "user"
	cfg.Contexts["context1"].Cluster = "cluster1"
	cfg.Contexts["context1"].Namespace = "namespace1"
	cfg.Contexts["context2"].AuthInfo = "user"
	cfg.Contexts["context2"].Cluster = "cluster2"
	cfg.Contexts["context2"].Namespace = "namespace2"
	cfg.CurrentContext = "context1"
	cfgData, err := clientcmd.Write(*cfg)
	assert.Check(t, err)
	_, err = kcFile.Write(cfgData)
	assert.Check(t, err)
	kcFile.Close()

	assert.Check(t, SetKubenetesContextKubeconfig(store, "external-default-context", kcFile.Name(), "", "", false))
	assert.Check(t, SetKubenetesContextKubeconfig(store, "external-context2", kcFile.Name(), "context2", "namespace-override", false))
	assert.Check(t, SetKubenetesContextKubeconfig(store, "embed-default-context", kcFile.Name(), "", "", true))
	assert.Check(t, SetKubenetesContextKubeconfig(store, "embed-context2", kcFile.Name(), "context2", "namespace-override", true))

	rawNotlsMeta, err := store.GetContextMetadata("raw-notls")
	assert.Check(t, err)
	rawNotlsSkipMeta, err := store.GetContextMetadata("raw-notls-skip")
	assert.Check(t, err)
	rawTlsMeta, err := store.GetContextMetadata("raw-tls")
	assert.Check(t, err)
	externalDefaultMeta, err := store.GetContextMetadata("external-default-context")
	assert.Check(t, err)
	externalContext2Meta, err := store.GetContextMetadata("external-context2")
	assert.Check(t, err)
	embededDefaultMeta, err := store.GetContextMetadata("embed-default-context")
	assert.Check(t, err)
	embededContext2Meta, err := store.GetContextMetadata("embed-context2")
	assert.Check(t, err)

	rawNoTls := ParseContext("raw-notls", rawNotlsMeta)
	rawNotlsSkip := ParseContext("raw-notls-skip", rawNotlsSkipMeta)
	rawTls := ParseContext("raw-tls", rawTlsMeta)
	externalDefault := ParseContext("external-default-context", externalDefaultMeta)
	externalContext2 := ParseContext("external-context2", externalContext2Meta)
	embededDefault := ParseContext("embed-default-context", embededDefaultMeta)
	embededContext2 := ParseContext("embed-context2", embededContext2Meta)

	rawNoTlsConfig, err := rawNoTls.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, rawNoTlsConfig, "https://test", "test", nil, nil, nil, false)
	rawNotlsSkipConfig, err := rawNotlsSkip.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, rawNotlsSkipConfig, "https://test", "test", nil, nil, nil, true)
	rawTlsConfig, err := rawTls.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, rawTlsConfig, "https://test", "test", []byte("ca"), []byte("cert"), []byte("key"), true)
	externalDefaultConfig, err := externalDefault.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, externalDefaultConfig, "https://server1", "namespace1", nil, []byte("cert"), []byte("key"), true)
	externalContext2Config, err := externalContext2.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, externalContext2Config, "https://server2", "namespace-override", []byte("ca"), []byte("cert"), []byte("key"), false)
	embededDefaultConfig, err := embededDefault.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, embededDefaultConfig, "https://server1", "namespace1", nil, []byte("cert"), []byte("key"), true)
	embededContext2Config, err := embededContext2.LoadKubernetesConfig(store)
	assert.Check(t, err)
	checkClientConfig(t, embededContext2Config, "https://server2", "namespace-override", []byte("ca"), []byte("cert"), []byte("key"), false)
}

func checkClientConfig(t *testing.T, config clientcmd.ClientConfig, server, namespace string, ca, cert, key []byte, skipTLSVerify bool) {
	cfg, err := config.ClientConfig()
	assert.Check(t, err)
	ns, _, _ := config.Namespace()
	assert.Equal(t, server, cfg.Host)
	assert.Equal(t, namespace, ns)
	assert.DeepEqual(t, ca, cfg.CAData)
	assert.DeepEqual(t, cert, cfg.CertData)
	assert.DeepEqual(t, key, cfg.KeyData)
	assert.Equal(t, skipTLSVerify, cfg.Insecure)
}
