package kubernetes

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/cli/cli/context/common"
	"github.com/docker/cli/cli/context/store"
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func testEndpoint(name, server, defaultNamespace string, ca, cert, key []byte, skipTLSVerify bool) Endpoint {
	var tlsData *common.TLSData
	if ca != nil || cert != nil || key != nil {
		tlsData = &common.TLSData{
			CA:   ca,
			Cert: cert,
			Key:  key,
		}
	}
	return Endpoint{
		EndpointMeta: EndpointMeta{
			EndpointMeta: common.EndpointMeta{
				ContextName:   name,
				Host:          server,
				SkipTLSVerify: skipTLSVerify,
			},
			DefaultNamespace: defaultNamespace,
		},
		TLSData: tlsData,
	}
}

func TestSaveLoadContexts(t *testing.T) {
	storeDir, err := ioutil.TempDir("", "test-load-save-k8-context")
	assert.NilError(t, err)
	defer os.RemoveAll(storeDir)
	store, err := store.New(storeDir)
	assert.NilError(t, err)
	assert.NilError(t, Save(store, testEndpoint("raw-notls", "https://test", "test", nil, nil, nil, false)))
	assert.NilError(t, Save(store, testEndpoint("raw-notls-skip", "https://test", "test", nil, nil, nil, true)))
	assert.NilError(t, Save(store, testEndpoint("raw-tls", "https://test", "test", []byte("ca"), []byte("cert"), []byte("key"), true)))

	kcFile, err := ioutil.TempFile(os.TempDir(), "test-load-save-k8-context")
	assert.NilError(t, err)
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
	assert.NilError(t, err)
	_, err = kcFile.Write(cfgData)
	assert.NilError(t, err)
	kcFile.Close()

	epDefault, err := FromKubeConfig("embed-default-context", kcFile.Name(), "", "")
	assert.NilError(t, err)
	epContext2, err := FromKubeConfig("embed-context2", kcFile.Name(), "context2", "namespace-override")
	assert.NilError(t, err)
	assert.NilError(t, Save(store, epDefault))
	assert.NilError(t, Save(store, epContext2))

	rawNoTLSMeta, err := store.GetContextMetadata("raw-notls")
	assert.NilError(t, err)
	rawNoTLSSkipMeta, err := store.GetContextMetadata("raw-notls-skip")
	assert.NilError(t, err)
	rawTLSMeta, err := store.GetContextMetadata("raw-tls")
	assert.NilError(t, err)
	embededDefaultMeta, err := store.GetContextMetadata("embed-default-context")
	assert.NilError(t, err)
	embededContext2Meta, err := store.GetContextMetadata("embed-context2")
	assert.NilError(t, err)

	rawNoTLS := Parse("raw-notls", rawNoTLSMeta)
	rawNoTLSSkip := Parse("raw-notls-skip", rawNoTLSSkipMeta)
	rawTLS := Parse("raw-tls", rawTLSMeta)
	embededDefault := Parse("embed-default-context", embededDefaultMeta)
	embededContext2 := Parse("embed-context2", embededContext2Meta)

	rawNoTLSEP, err := rawNoTLS.WithTLSData(store)
	assert.NilError(t, err)
	checkClientConfig(t, store, rawNoTLSEP, "https://test", "test", nil, nil, nil, false)
	rawNoTLSSkipEP, err := rawNoTLSSkip.WithTLSData(store)
	assert.NilError(t, err)
	checkClientConfig(t, store, rawNoTLSSkipEP, "https://test", "test", nil, nil, nil, true)
	rawTLSEP, err := rawTLS.WithTLSData(store)
	assert.NilError(t, err)
	checkClientConfig(t, store, rawTLSEP, "https://test", "test", []byte("ca"), []byte("cert"), []byte("key"), true)
	embededDefaultEP, err := embededDefault.WithTLSData(store)
	assert.NilError(t, err)
	checkClientConfig(t, store, embededDefaultEP, "https://server1", "namespace1", nil, []byte("cert"), []byte("key"), true)
	embededContext2EP, err := embededContext2.WithTLSData(store)
	assert.NilError(t, err)
	checkClientConfig(t, store, embededContext2EP, "https://server2", "namespace-override", []byte("ca"), []byte("cert"), []byte("key"), false)
}

func checkClientConfig(t *testing.T, s store.Store, ep Endpoint, server, namespace string, ca, cert, key []byte, skipTLSVerify bool) {
	config, err := ep.KubernetesConfig()
	assert.NilError(t, err)
	cfg, err := config.ClientConfig()
	assert.NilError(t, err)
	ns, _, _ := config.Namespace()
	assert.Equal(t, server, cfg.Host)
	assert.Equal(t, namespace, ns)
	assert.DeepEqual(t, ca, cfg.CAData)
	assert.DeepEqual(t, cert, cfg.CertData)
	assert.DeepEqual(t, key, cfg.KeyData)
	assert.Equal(t, skipTLSVerify, cfg.Insecure)
}
