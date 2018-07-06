package context

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/docker/context-store"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/pkg/errors"
)

// Context is a typed wrapper around a context-store context
type Context struct {
	Name          string
	Host          string
	SkipTLSVerify bool
	APIVersion    string
}

// LoadTLSConfig extracts a context docker endpoint TLS config
func (c *Context) LoadTLSConfig(s store.Store) (*tls.Config, error) {
	tlsFiles, err := s.ListContextTLSFiles(c.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve context tls info: %s", err)
	}
	if epTLSFiles, ok := tlsFiles[dockerEndpointKey]; ok || c.SkipTLSVerify {
		var tlsOpts []func(*tls.Config)
		var keyBytes, certBytes []byte
		for _, f := range epTLSFiles {
			data, err := s.GetContextTLSData(c.Name, dockerEndpointKey, f)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to retrieve context tls info: %s", err)
			}
			switch f {
			case caKey:
				certPool := x509.NewCertPool()
				if !certPool.AppendCertsFromPEM(data) {
					return nil, errors.New("failed to retrieve context tls info: ca.pem seems invalid")
				}
				tlsOpts = append(tlsOpts, func(cfg *tls.Config) {
					cfg.RootCAs = certPool
				})
			case certKey:
				certBytes = data
			case keyKey:
				keyBytes = data
			}
		}
		if keyBytes != nil && certBytes != nil {
			x509cert, err := tls.X509KeyPair(certBytes, keyBytes)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to retrieve context tls info: %s", err)
			}
			tlsOpts = append(tlsOpts, func(cfg *tls.Config) {
				cfg.Certificates = []tls.Certificate{x509cert}
			})
		}
		if c.SkipTLSVerify {
			tlsOpts = append(tlsOpts, func(cfg *tls.Config) {
				cfg.InsecureSkipVerify = true
			})
		}
		return tlsconfig.ClientDefault(tlsOpts...), nil
	}
	return nil, nil
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

// Parse parses a context docker endpoint metadata into a typed Context structure
func Parse(name string, metadata store.ContextMetadata) (*Context, error) {
	ep, ok := metadata.Endpoints[dockerEndpointKey]
	if !ok {
		return nil, errors.New("cannot find docker endpoint in context")
	}
	host, _ := getMetaString(ep, hostKey)
	skipTLSVerify, _ := getMetaBool(ep, skipTLSVerifyKey)
	apiVersion, _ := getMetaString(ep, apiVersionKey)
	return &Context{
		Name:          name,
		Host:          host,
		SkipTLSVerify: skipTLSVerify,
		APIVersion:    apiVersion,
	}, nil
}
