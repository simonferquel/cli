package docker

import (
	"os"

	"github.com/docker/cli/cli/context/store"
)

// ToStoreMeta convert the endpoint to the store format
func (e *EndpointMeta) ToStoreMeta() store.Metadata {
	meta := e.EndpointMeta.ToStoreMeta()
	if e.APIVersion != "" {
		meta[apiVersionKey] = e.APIVersion
	}
	return meta
}

// Save saves the docker endpoint in a context store
func Save(s store.Store, endpoint Endpoint) error {
	ctxMeta, err := s.GetContextMetadata(endpoint.ContextName)
	switch {
	case os.IsNotExist(err):
		ctxMeta = store.ContextMetadata{
			Endpoints: make(map[string]store.Metadata),
			Metadata:  make(store.Metadata),
		}
	case err != nil:
		return err
	}
	ctxMeta.Endpoints[dockerEndpointKey] = endpoint.ToStoreMeta()
	if err := s.CreateOrUpdateContext(endpoint.ContextName, ctxMeta); err != nil {
		return err
	}
	return s.ResetContextEndpointTLSMaterial(endpoint.ContextName, dockerEndpointKey, endpoint.TLSData.ToStoreTLSData())
}
