package common

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

var Storage StorageGateway

type StorageGateway interface {
	GetAllKeys(path string) ([]string, error)
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Exists(key string) (bool, error)
	Delete(key string) error
}

func ExtractID(key string) string {
	p := filepath.Base(key)
	if strings.HasSuffix(p, ".pbf") {
		return p[:len(p)-4]
	}
	return p
}

func IDToFile(ID string) string {
	if strings.HasSuffix(ID, ".pbf") {
		return ID
	}
	return ID + ".pbf"
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
	return lfs.s.GetAllKeys(path)
}

func (lfs *ReadOnlyStorage) Get(key string) ([]byte, error) {
	return lfs.s.Get(key)
}

func (lfs *ReadOnlyStorage) Put(key string, data []byte) error {
	return nil
}

func (lfs *ReadOnlyStorage) Exists(key string) (bool, error) {
	return lfs.s.Exists(key)
}

func (lfs *ReadOnlyStorage) Delete(key string) error {
	return nil
}

func LoadFromStorage[T proto.Message](path string, item T) error {
	data, err := Storage.Get(path)
	if err != nil {
		log.Printf("Error loading character  %v : %v", path, err)
		return nil
	}
	err = proto.Unmarshal(data, item)
	if err != nil {
		log.Printf("Error parsing %v : %v", path, err)
		return nil
	}
	return nil
}

func SaveToStorage[T proto.Message](path string, item T) error {
	data, err := proto.Marshal(item)
	if err != nil {
		log.Printf("Error marshalling  %v : %v", path, err)
		return nil
	}
	err = Storage.Put(path, data)
	if err != nil {
		log.Printf("Error saving %v : %v", path, err)
		return nil
	}
	return nil
}

func LoadAllFromStorage[T proto.Message](path string, factory func() T) ([]T, []string, error) {
	keys, err := Storage.GetAllKeys(path)
	if err != nil {
		log.Printf("Error Listing  %v : %v", path, err)
		return nil, nil, nil
	}

	var rtn []T
	var rtnIds []string
	for _, key := range keys {
		item := factory()
		keypath := filepath.Join(path, key)
		err = LoadFromStorage(keypath, item)
		if err == nil {
			rtn = append(rtn, item)
			rtnIds = append(rtnIds, ExtractID(key))
		}
	}
	return rtn, rtnIds, nil
}
