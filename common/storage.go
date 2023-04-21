package common

import (
	"errors"
	"os"
	"path/filepath"
)

var Storage StorageGateway

type StorageGateway interface {
	GetAllKeys(path string) ([]string, error)
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Exists(key string) (bool, error)
	Delete(key string) error
}

type LocalFSStorage struct {
	rootDir string
}

func (lfs *LocalFSStorage) GetAllKeys(path string) ([]string, error) {
	full := filepath.Join(lfs.rootDir, path)
	entries, err := os.ReadDir(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var rtn []string
	for _, e := range entries {
		if !e.IsDir() {
			rtn = append(rtn, e.Name())
		}
	}
	return rtn, nil
}

func (lfs *LocalFSStorage) Get(key string) ([]byte, error) {
	path := filepath.Join(lfs.rootDir, key)
	return os.ReadFile(path)
}

func (lfs *LocalFSStorage) Put(key string, data []byte) error {
	path := filepath.Join(lfs.rootDir, key)
	dirs := filepath.Dir(path)
	err := os.MkdirAll(dirs, 0600)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (lfs *LocalFSStorage) Exists(key string) (bool, error) {
	path := filepath.Join(lfs.rootDir, key)
	_, err := os.Stat(path)
	// perr := err.(*os.PathError)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lfs *LocalFSStorage) Delete(key string) error {
	path := filepath.Join(lfs.rootDir, key)
	return os.Remove(path)
}

type ReadOnlyStorage struct {
	s StorageGateway
}

func (lfs *ReadOnlyStorage) GetAllKeys(path string) ([]string, error) {
	return lfs.GetAllKeys(path)
}

func (lfs *ReadOnlyStorage) Get(key string) ([]byte, error) {
	return lfs.Get(key)
}

func (lfs *ReadOnlyStorage) Put(key string, data []byte) error {
	return nil
}

func (lfs *ReadOnlyStorage) Exists(key string) (bool, error) {
	return lfs.Exists(key)
}

func (lfs *ReadOnlyStorage) Delete(key string) error {
	return nil
}
