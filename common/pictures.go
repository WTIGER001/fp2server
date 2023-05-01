package common

import (
	"path/filepath"

	"google.golang.org/protobuf/proto"
)

var Pictures = newPictureManager()

type PictureManager struct{}

func newPictureManager() *PictureManager {
	return &PictureManager{}
}

func (p *PictureManager) toGameId(gameID string, modelType ModelType) string {
	if int32(modelType) >= 50 {
		return "_references"
	}
	if gameID == "" {
		if ActiveGame != nil {
			return ActiveGame.ID
		}
	}
	return gameID
}

func (p *PictureManager) Key(id string, modelType ModelType, gameId string, tag string) string {
	gameId = p.toGameId(gameId, modelType)
	if tag == "" {
		tag = "picture"
	}
	mtName := ModelType_name[int32(modelType)]
	return filepath.Join("pictures", gameId, mtName, id, tag+".pbf")
}

func (p *PictureManager) Get(id string, modelType ModelType, gameId string, tag string) (*Picture, error) {
	key := p.Key(id, modelType, gameId, tag)
	data, err := Storage.Get(key)
	if err != nil {
		return nil, err
	}

	pic := &Picture{}
	err = proto.Unmarshal(data, pic)
	if err != nil {
		return nil, err
	}
	return pic, nil
}

func (p *PictureManager) Set(pic *Picture) error {
	key := p.Key(pic.ID, pic.Type, pic.GameID, pic.Tag)

	data, err := proto.Marshal(pic)
	if err != nil {
		return err
	}
	return Storage.Put(key, data)
}

func (p *PictureManager) Delete(pic *Picture) error {
	key := p.Key(pic.ID, pic.Type, pic.GameID, pic.Tag)

	return Storage.Delete(key)
}
