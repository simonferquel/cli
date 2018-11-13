package docker

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/cli/cli/context/common"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/pkg/errors"
)

// EndpointMeta is a typed wrapper around a context-store generic endpoint describing
// a Docker Engine endpoint, without its tls config
type EndpointMeta struct {
	common.EndpointMeta
	APIVersion string
}

// Endpoint is a typed wrapper around a context-store generic endpoint describing
// a Docker Engine endpoint, with its tls data
type Endpoint struct {
	EndpointMeta
	TLSData     *common.TLSData
	TLSPassword string
}

// WithTLSData loads TLS materials for the endpoint
func (c *EndpointMeta) WithTLSData(s store.Store) (Endpoint, error) {
	tlsData, err := common.LoadTLSData(s, c.ContextName, dockerEndpointKey)
	if err != nil {
		return Endpoint{}, err
	}
	return Endpoint{
		EndpointMeta: *c,
		TLSData:      tlsData,
	}, nil
}

// tlsConfig extracts a context docker endpoint TLS config
func (c *Endpoint) tlsConfig() (*tls.Config, error) {
	if c.TLSData == nil && !c.SkipTLSVerify {
		// there is no specific tls config
		return nil, nil
	}
	var tlsOpts []func(*tls.Config)
	if c.TLSData != nil && c.TLSData.CA != nil {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(c.TLSData.CA) {
			return nil, errors.New("failed to retrieve context tls info: ca.pem seems invalid")
		}
		tlsOpts = append(tlsOpts, func(cfg *tls.Config) {
			cfg.RootCAs = certPool
		})
	}
	if c.TLSData != nil && c.TLSData.Key != nil && c.TLSData.Cert != nil {
		keyBytes := c.TLSData.Key
		pemBlock, _ := pem.Decode(keyBytes)
		if pemBlock == nil {
			return nil, fmt.Errorf("no valid private key found")
		}

		var err error
		if x509.IsEncryptedPEMBlock(pemBlock) {
			keyBytes, err = x509.DecryptPEMBlock(pemBlock, []byte(c.TLSPassword))
			if err != nil {
				return nil, errors.Wrap(err, "private key is encrypted, but could not decrypt it")
			}
			keyBytes = pem.EncodeToMemory(&pem.Block{Type: pemBlock.Type, Bytes: keyBytes})
		}

		x509cert, err := tls.X509KeyPair(c.TLSData.Cert, keyBytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve context tls info")
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

// ConfigureClient configures a docker client
func (c *Endpoint) ConfigureClient(cli *client.Client) error {
	if c.Host != "" {
		helper, err := connhelper.GetConnectionHelper(c.Host)
		if err != nil {
			return err
		}
		if helper == nil {
			if err := client.WithHost(c.Host)(cli); err != nil {
				return err
			}
		} else {
			httpClient := &http.Client{
				// No tls
				// No proxy
				Transport: &http.Transport{
					DialContext: helper.Dialer,
				},
			}
			if err := client.WithHTTPClient(httpClient)(cli); err != nil {
				return err
			}
			if err := client.WithHost(helper.Host)(cli); err != nil {
				return err
			}
			if err := client.WithDialContext(helper.Dialer)(cli); err != nil {
				return err
			}
		}
	}
	tlsConfig, err := c.tlsConfig()
	if err != nil {
		return err
	}
	if tlsConfig != nil {
		httpClient := cli.HTTPClient()
		if transport, ok := httpClient.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = tlsConfig
		} else {
			return errors.Errorf("cannot apply tls config to transport: %T", httpClient.Transport)
		}
		if err := client.WithHTTPClient(httpClient)(cli); err != nil {
			return err
		}
	}
	version := os.Getenv("DOCKER_API_VERSION")
	if version == "" {
		version = c.APIVersion
	}
	if version != "" {
		if err := client.WithVersion(version)(cli); err != nil {
			return err
		}
	}
	return nil
}

// Parse parses a context docker endpoint metadata into a typed EndpointMeta structure
func Parse(name string, metadata store.ContextMetadata) (EndpointMeta, error) {
	ep, ok := metadata.Endpoints[dockerEndpointKey]
	if !ok {
		return EndpointMeta{}, errors.New("cannot find docker endpoint in context")
	}
	commonMeta := common.Parse(name, dockerEndpointKey, metadata)
	if commonMeta == nil {
		return EndpointMeta{}, errors.New("cannot find docker endpoint in context")
	}
	apiVersion, _ := ep.GetString(apiVersionKey)
	return EndpointMeta{
		EndpointMeta: *commonMeta,
		APIVersion:   apiVersion,
	}, nil
}
