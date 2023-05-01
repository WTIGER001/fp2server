package main

import (
	"github.com/wtiger001/fp2server/common"
)

type CrudMessageHandler struct {
}

func (mh *CrudMessageHandler) Handle(m *common.Fp2Message) {
	switch m.Data.(type) {

	case *common.Fp2Message_GetRequest:
		mh.OnGet(m)
	case *common.Fp2Message_GetAllRequest:
		mh.OnGetAll(m)
	case *common.Fp2Message_ListRequest:

	case *common.Fp2Message_UpdateRequest:
		mh.OnUpdate(m)
	}

}

func ModelGet(modelType common.ModelType, id string) (interface{}, error) {
	switch modelType {
	case common.ModelType_ModelType_Armor:
		return nil, nil
	case common.ModelType_ModelType_Character:
		if common.ActiveGame != nil {
			return common.ActiveGame.Characters().Get(id)
		}
		return nil, nil
	case common.ModelType_ModelType_Game:
		return common.Games.Get(id)
	case common.ModelType_ModelType_Orb:
		return nil, nil
	case common.ModelType_ModelType_Player:
		return nil, nil
	case common.ModelType_ModelType_RefArmor:
		return common.References.Armors.Get(id)
	case common.ModelType_ModelType_RefGameTerm:
		return common.References.GameTerms.Get(id)
	case common.ModelType_ModelType_RefOrb:
		return common.References.Orbs.Get(id)
	case common.ModelType_ModelType_RefSkill:
		return common.References.Skills.Get(id)
	case common.ModelType_ModelType_RefWeapon:
		return common.References.Weapons.Get(id)
	case common.ModelType_ModelType_Weapon:
		return nil, nil
	case common.ModelType_ModelType_Unkown:
		return nil, nil
	}
	return nil, nil
}

func ModelGetAll(modelType common.ModelType) ([]*common.Model, error) {
	switch modelType {
	case common.ModelType_ModelType_Armor:
		return nil, nil
	case common.ModelType_ModelType_Character:
		if common.ActiveGame != nil {
			items, err := common.ActiveGame.Characters().GetAll()
			if err != nil {
				return nil, err
			}
			models := common.ToModels(items)
			return models, err
		}
		return nil, nil
	case common.ModelType_ModelType_Game:
		items, err := common.Games.GetAll()
		if err != nil {
			return nil, err
		}
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_Orb:
		return nil, nil
	case common.ModelType_ModelType_Player:
		return nil, nil
	case common.ModelType_ModelType_RefArmor:
		items, err := common.References.Armors.GetAll()
		if err != nil {
			return nil, err
		}
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_RefGameTerm:
		items, err := common.References.GameTerms.GetAll()
		if err != nil {
			return nil, err
		}
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_RefOrb:
		items, err := common.References.Orbs.GetAll()
		if err != nil {
			return nil, err
		}
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_RefSkill:
		items, err := common.References.Skills.GetAll()
		if err != nil {
			return nil, err
		}
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_RefWeapon:
		items, err := common.References.Weapons.GetAll()
		if err != nil {
			return nil, err
		}
		return common.ToModels(items), err
		models := common.ToModels(items)
		return models, err
	case common.ModelType_ModelType_Weapon:
		return nil, nil
	case common.ModelType_ModelType_Unkown:
		return nil, nil
	}
	return nil, nil
}

func ModelUpdate(modelType common.ModelType, id string) (interface{}, error) {
	switch modelType {
	case common.ModelType_ModelType_Armor:
		return nil, nil
	case common.ModelType_ModelType_Character:
		if common.ActiveGame != nil {
			return common.ActiveGame.Characters().Get(id)
		}
		return nil, nil
	case common.ModelType_ModelType_Game:
		return common.Games.Get(id)
	case common.ModelType_ModelType_Orb:
		return nil, nil
	case common.ModelType_ModelType_Player:
		return nil, nil

	case common.ModelType_ModelType_Picture:

	case common.ModelType_ModelType_RefArmor:
		return common.References.Armors.Get(id)
	case common.ModelType_ModelType_RefGameTerm:
		return common.References.GameTerms.Get(id)
	case common.ModelType_ModelType_RefOrb:
		return common.References.Orbs.Get(id)
	case common.ModelType_ModelType_RefSkill:
		return common.References.Skills.Get(id)
	case common.ModelType_ModelType_RefWeapon:
		return common.References.Weapons.Get(id)
	case common.ModelType_ModelType_Weapon:
		return nil, nil
	case common.ModelType_ModelType_Unkown:
		return nil, nil
	}
	return nil, nil
}

func ModelList(modelType common.ModelType) ([]*common.Model, error) {
	return nil, nil
}

func (mh *CrudMessageHandler) OnGet(m *common.Fp2Message) {
	request := m.GetGetRequest()
	item, err := ModelGet(request.Type, request.ID)
	if err != nil {
		SendError(err, m)
		return
	}

	mp.SendTo(m.Sender, &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Data: &common.Fp2Message_GetResponse{
			GetResponse: &common.GetResponse{
				Type:  request.Type,
				Model: common.ToModel(item),
			},
		},
	})
}

func (mh *CrudMessageHandler) OnGetAll(m *common.Fp2Message) {
	request := m.GetGetAllRequest()
	items, err := ModelGetAll(request.Type)
	if err != nil {
		SendError(err, m)
		return
	}

	mp.SendTo(m.Sender, &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Data: &common.Fp2Message_GetAllResponse{
			GetAllResponse: &common.GetAllResponse{
				Type:  request.Type,
				Items: items,
			},
		},
	})
}
func (mh *CrudMessageHandler) OnUpdate(m *common.Fp2Message) {
	var modelRtn *common.Model
	var err error
	request := m.GetUpdateRequest()
	switch request.Type {
	case common.UpdateType_UT_Save:
		modelRtn, err = mh.Update(request)
	case common.UpdateType_UT_Delete:
		err = mh.Delete(request)
	}

	if err != nil {
		SendError(err, m)
		return
	}

	// Send Response
	mp.Broadcast(&common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Data: &common.Fp2Message_ModelChangedEvent{
			ModelChangedEvent: &common.ModelChangedEvent{
				Type:         request.Type,
				UpdateReason: request.UpdateReason,
				Model:        modelRtn,
			},
		},
	})
}

func (mh *CrudMessageHandler) Update(r *common.UpdateRequest) (*common.Model, error) {
	switch r.Model.Data.(type) {
	case *common.Model_Armor:
		// Do nothing
	case *common.Model_Character:
		if common.ActiveGame != nil {
			item := r.Model.GetCharacter()
			common.ActiveGame.Characters().Set(item, item.ID)
		}
	case *common.Model_Game:
		item := r.Model.GetGame()
		err := common.Games.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_Orb:
		// Do Nothing
	case *common.Model_Player:
		// Do Nothing
	case *common.Model_Picture:
		item := r.Model.GetPicture()
		err := common.Pictures.Set(item)
		return nil, err
	case *common.Model_RefArmor:
		item := r.Model.GetRefArmor()
		err := common.References.Armors.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_RefGameTerm:
		item := r.Model.GetRefGameTerm()
		err := common.References.GameTerms.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_RefOrb:
		item := r.Model.GetRefOrb()
		err := common.References.Orbs.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_RefSkill:
		item := r.Model.GetRefSkill()
		err := common.References.Skills.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_RefWeapon:
		item := r.Model.GetRefWeapon()
		err := common.References.Weapons.Set(item, item.ID)
		return common.ToModel(item), err
	case *common.Model_Weapon:
		// Do Nothing
	}
	return nil, nil
}

func (mh *CrudMessageHandler) Delete(r *common.UpdateRequest) error {
	switch r.Model.Data.(type) {
	case *common.Model_Armor:
		// Do nothing
	case *common.Model_Character:
		if common.ActiveGame != nil {
			item := r.Model.GetCharacter()
			return common.ActiveGame.Characters().Delete(item.ID)
		}
	case *common.Model_Game:
		item := r.Model.GetGame()
		return common.Games.Delete(item.ID)
	case *common.Model_Orb:
		// Do Nothing
	case *common.Model_Player:
		// Do Nothing
	case *common.Model_Picture:
		item := r.Model.GetPicture()
		common.Pictures.Delete(item)
	case *common.Model_RefArmor:
		item := r.Model.GetRefArmor()
		return common.References.Armors.Delete(item.ID)
	case *common.Model_RefGameTerm:
		item := r.Model.GetRefGameTerm()
		return common.References.GameTerms.Delete(item.ID)
	case *common.Model_RefOrb:
		item := r.Model.GetRefOrb()
		return common.References.Orbs.Delete(item.ID)
	case *common.Model_RefSkill:
		item := r.Model.GetRefSkill()
		return common.References.Skills.Delete(item.ID)
	case *common.Model_RefWeapon:
		item := r.Model.GetRefWeapon()
		return common.References.Weapons.Delete(item.ID)
	case *common.Model_Weapon:
		// Do Nothing
	}
	return nil
}

func SendError(err error, m *common.Fp2Message) {
	mp.SendTo(m.Sender, &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Data: &common.Fp2Message_ErrorResponse{
			&common.ErrorResponse{
				Error:          err.Error(),
				ErrorCode:      0,
				ErroredMessage: m,
			},
		},
	})
}
