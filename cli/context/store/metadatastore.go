package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	metadataDir = "meta"
	metaFile    = "meta.json"
)

type metadataStore struct {
	root string
}

func (s *metadataStore) contextDir(name string) string {
	return filepath.Join(s.root, name)
}

func (s *metadataStore) createOrUpdate(name string, meta ContextMetadata) error {
	contextDir := s.contextDir(name)
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return err
	}
	bytes, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(contextDir, metaFile), bytes, 0644)
}

func (s *metadataStore) get(name string) (ContextMetadata, error) {
	contextDir := s.contextDir(name)
	bytes, err := ioutil.ReadFile(filepath.Join(contextDir, metaFile))
	if err != nil {
		return ContextMetadata{}, convertContextDoesNotExist(name, err)
	}
	var r ContextMetadata
	err = json.Unmarshal(bytes, &r)
	return r, err
}

func (s *metadataStore) remove(name string) error {
	contextDir := s.contextDir(name)
	return os.RemoveAll(contextDir)
}

func (s *metadataStore) list() (map[string]ContextMetadata, error) {
	ctxNames, err := listRecursivelyMetadataDirs(s.root)
	if err != nil {
		if os.IsNotExist(err) {
			// store is empty, meta dir does not exist yet
			// this should not be considered an error
			return map[string]ContextMetadata{}, nil
		}
		return nil, err
	}
	res := make(map[string]ContextMetadata)
	for _, name := range ctxNames {
		res[name], err = s.get(name)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func isContextDir(path string) bool {
	s, err := os.Stat(filepath.Join(path, metaFile))
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func listRecursivelyMetadataDirs(root string) ([]string, error) {
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, fi := range fis {
		if fi.IsDir() {
			if isContextDir(filepath.Join(root, fi.Name())) {
				result = append(result, fi.Name())
			}
			subs, err := listRecursivelyMetadataDirs(filepath.Join(root, fi.Name()))
			if err != nil {
				return nil, err
			}
			for _, s := range subs {
				result = append(result, fmt.Sprintf("%s/%s", fi.Name(), s))
			}
		}
	}
	return result, nil
}

func convertContextDoesNotExist(name string, err error) error {
	if os.IsNotExist(err) {
		return &contextDoesNotExistError{name: name}
	}
	return err
}
