package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

// Store provides a context store for easily remembering endpoints configuration
type Store interface {
	GetCurrentContext() string
	SetCurrentContext(name string) error
	ListContexts() (map[string]ContextMetadata, error)
	CreateOrUpdateContext(name string, meta ContextMetadata) error
	GetContextMetadata(name string) (ContextMetadata, error)
	ResetContextTLSMaterial(name string, data *ContextTLSData) error
	ResetContextEndpointTLSMaterial(contextName string, endpointName string, data *EndpointTLSData) error
	ListContextTLSFiles(name string) (map[string]EndpointFiles, error)
	GetContextTLSData(contextName, endpointName, fileName string) ([]byte, error)
}
type store struct {
	configFile     string
	currentContext string
	meta           *metadataStore
	tls            *tlsStore
}

// NewStore creates a store from a given directory.
// If the directory does not exist or is empty, initialize it
func NewStore(dir string) (Store, error) {
	metaRoot := filepath.Join(dir, metadataDir)
	tlsRoot := filepath.Join(dir, tlsDir)
	configFile := filepath.Join(dir, configFileName)
	err := os.MkdirAll(metaRoot, 0755)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(tlsRoot, 0700)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(configFile)

	switch {
	case os.IsNotExist(err):
		//create default file
		err = ioutil.WriteFile(configFile, []byte("{}"), 0644)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg config
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &store{
		configFile:     configFile,
		currentContext: cfg.CurrentContext,
		meta: &metadataStore{
			root: metaRoot,
		},
		tls: &tlsStore{
			root: tlsRoot,
		},
	}, nil
}

func (s *store) GetCurrentContext() string {
	return s.currentContext
}

func (s *store) SetCurrentContext(name string) error {
	cfg := config{CurrentContext: name}
	configBytes, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.configFile, configBytes, 0644)
	if err != nil {
		return err
	}
	s.currentContext = name
	return nil
}

func (s *store) ListContexts() (map[string]ContextMetadata, error) {
	return s.meta.list()
}

func (s *store) CreateOrUpdateContext(name string, meta ContextMetadata) error {
	return s.meta.createOrUpdate(name, meta)
}

func (s *store) GetContextMetadata(name string) (ContextMetadata, error) {
	return s.meta.get(name)
}

func (s *store) ResetContextTLSMaterial(name string, data *ContextTLSData) error {
	err := s.tls.removeAllContextData(name)
	if err != nil {
		return err
	}
	if data != nil {
		for ep, files := range data.Endpoints {
			for fileName, data := range files.Files {
				err = s.tls.createOrUpdate(name, ep, fileName, data)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *store) ResetContextEndpointTLSMaterial(contextName string, endpointName string, data *EndpointTLSData) error {
	err := s.tls.removeAllEndpointData(contextName, endpointName)
	if err != nil {
		return err
	}
	if data != nil {
		for fileName, data := range data.Files {
			err = s.tls.createOrUpdate(contextName, endpointName, fileName, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *store) ListContextTLSFiles(name string) (map[string]EndpointFiles, error) {
	return s.tls.listContextData(name)
}

func (s *store) GetContextTLSData(contextName, endpointName, fileName string) ([]byte, error) {
	return s.tls.getData(contextName, endpointName, fileName)
}

type config struct {
	CurrentContext string `json:"current_context,omitempty"`
}

// EndpointTLSData represents tls data for a given endpoint
type EndpointTLSData struct {
	Files map[string][]byte
}

// ContextTLSData represents tls data for a whole context
type ContextTLSData struct {
	Endpoints map[string]EndpointTLSData
}
