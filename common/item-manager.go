package common

import (
	"path/filepath"

	"google.golang.org/protobuf/proto"
)

func NewItemManager[T proto.Message](path string, factory func() T) *ItemManager[T] {
	return &ItemManager[T]{
		factory: factory,
		path:    path,
		cache:   NewCache[T](),
	}
}

type ItemManager[T proto.Message] struct {
	factory func() T
	path    string
	cache   Cache[T]
	onSave  func(item T) error
}

func (im *ItemManager[T]) Count() int {
	items, _ := im.GetAll()
	return len(items)
}

func (im *ItemManager[T]) Set(item T, id string) error {
	im.cache.Set(ExtractID(id), item)
	return im.Save(item, id)
}

func (im *ItemManager[T]) Get(idOrKey string) (T, error) {
	id := ExtractID(idOrKey)
	item, found := im.cache.Get(id)
	if found {
		return item, nil
	}
	item, err := im.Load(id)
	if err != nil {

	}
	im.cache.Set(id, item)
	return item, err
}

func (im *ItemManager[T]) GetAll() ([]T, error) {
	if !im.cache.IsEmpty() {
		return im.cache.All(), nil
	}
	items, ids, err := im.LoadAll()
	if err != nil {
		return nil, err
	}
	for i, item := range items {
		im.cache.Set(ids[i], item)
	}
	return items, nil
}

func (im *ItemManager[T]) Load(idOrKey string) (T, error) {
	key := IDToFile(idOrKey)
	path := filepath.Join(im.path, key)
	item := im.factory()
	err := LoadFromStorage(path, item)
	return item, err
}

func (im *ItemManager[T]) LoadAll() ([]T, []string, error) {
	rtn, ids, err := LoadAllFromStorage(im.path, im.factory)
	return rtn, ids, err
}

func (im *ItemManager[T]) Save(item T, id string) error {
	if im.onSave != nil {
		err := im.onSave(item)
		if err != nil {
			return err
		}
	}

	key := IDToFile(id)
	path := filepath.Join(im.path, key)
	return SaveToStorage(path, item)
}

func (im *ItemManager[T]) Delete(id string) error {
	key := IDToFile(id)
	path := filepath.Join(im.path, key)
	return Storage.Delete(path)
}

func (im *ItemManager[T]) Exists(id string) (bool, error) {
	key := IDToFile(id)
	path := filepath.Join(im.path, key)
	return Storage.Exists(path)
}
