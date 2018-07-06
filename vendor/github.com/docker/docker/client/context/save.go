package context

import (
	"os"

	"github.com/docker/context-store"
)

// SetDockerEndpoint set the docker endpoint of a context
func SetDockerEndpoint(s store.Store, name, host, apiVersion string, ca, cert, key []byte, skipTLSVerify bool) error {
	ctxMeta, err := s.GetContextMetadata(name)
	switch {
	case os.IsNotExist(err):
		ctxMeta = store.ContextMetadata{
			Endpoints: make(map[string]store.EndpointMetadata),
			Metadata:  make(map[string]interface{}),
		}
	case err != nil:
		return err
	}
	epMeta := make(store.EndpointMetadata)
	epMeta[hostKey] = host
	if apiVersion != "" {
		epMeta[apiVersionKey] = apiVersion
	}
	if skipTLSVerify {
		epMeta[skipTLSVerifyKey] = true
	}
	ctxMeta.Endpoints[dockerEndpointKey] = epMeta
	err = s.CreateOrUpdateContext(name, ctxMeta)
	if err != nil {
		return err
	}
	return s.ResetContextEndpointTLSMaterial(name, dockerEndpointKey, createEnpointTLSData(ca, cert, key))
}

func createEnpointTLSData(ca, cert, key []byte) *store.EndpointTLSData {
	if ca == nil && cert == nil && key == nil {
		return nil
	}
	result := store.EndpointTLSData{
		Files: make(map[string][]byte),
	}
	if ca != nil {
		result.Files[caKey] = ca
	}
	if cert != nil {
		result.Files[certKey] = cert
	}
	if key != nil {
		result.Files[keyKey] = key
	}
	return &result
}
