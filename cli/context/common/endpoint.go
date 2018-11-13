package common

import (
	"github.com/docker/cli/cli/context/store"
)

const (
	hostKey          = "host"
	skipTLSVerifyKey = "skipTLSVerify"
	caKey            = "ca.pem"
	certKey          = "cert.pem"
	keyKey           = "key.pem"
)

// EndpointMeta contains fields we expect to be common for all context endpoints
type EndpointMeta struct {
	ContextName   string
	Host          string
	SkipTLSVerify bool
}

// ToStoreMeta converts the endpoint to the store format
func (e *EndpointMeta) ToStoreMeta() store.Metadata {
	return store.Metadata{
		hostKey:          e.Host,
		skipTLSVerifyKey: e.SkipTLSVerify,
	}
}

// Parse parses a context endpoint metadata into a typed EndpointMeta structure
func Parse(contextName, endpointName string, metadata store.ContextMetadata) *EndpointMeta {
	ep, ok := metadata.Endpoints[endpointName]
	if !ok {
		return nil
	}
	host, _ := ep.GetString(hostKey)
	skipTLSVerify, _ := ep.GetBoolean(skipTLSVerifyKey)
	return &EndpointMeta{
		ContextName:   contextName,
		Host:          host,
		SkipTLSVerify: skipTLSVerify,
	}
}
