package main

import (
	"fmt"

	"github.com/wtiger001/fp2server/common"
)

func SendGetAllRequest(t common.ModelType, text string) *common.Fp2Message {
	response := MessageBus.SendAndRecieve(&common.Fp2Message{
		Data: &common.Fp2Message_GetAllRequest{
			GetAllRequest: &common.GetAllRequest{
				Type: common.ModelType_ModelType_RefWeapon,
			},
		},
	})
	return response
}

func ListModels(t common.ModelType, text string) ([]*common.IDName, error) {
	response := MessageBus.SendAndRecieve(&common.Fp2Message{
		Data: &common.Fp2Message_ListRequest{
			ListRequest: &common.ListRequest{
				Type: common.ModelType_ModelType_RefWeapon,
			},
		},
	})
	err := response.GetErrorResponse()
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error)
	}
	return response.GetListResponse().Items, nil
}

func GetAllWeaponRefs() []*common.RefWeapon {
	response := SendGetAllRequest(common.ModelType_ModelType_RefWeapon, "")

	// What to do on error
	// err := response.GetErrorResponse()

	ref := response.GetGetAllResponse()
	if ref == nil {
		return nil
	}
	items := ref.Items
	var rtn []*common.RefWeapon
	for _, item := range items {
		rtn = append(rtn, item.GetRefWeapon())
	}
	return rtn
}

func ModelGetAll(modelType common.ModelType) ([]*common.Model, error) {
	response := SendGetAllRequest(modelType, "")

	err := response.GetErrorResponse()
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error)
	}
	ref := response.GetGetAllResponse()
	return ref.Items, nil
}

func sendGetRequest(mt common.ModelType, id string) *common.Fp2Message {
	response := MessageBus.SendAndRecieve(&common.Fp2Message{
		Data: &common.Fp2Message_GetRequest{
			GetRequest: &common.GetRequest{
				Type: mt,
				ID:   id,
			},
		},
	})
	return response
}

func GetWeaponRef(ID string) *common.RefWeapon {
	response := sendGetRequest(common.ModelType_ModelType_RefWeapon, ID)
	ref := response.GetGetResponse()
	return ref.GetModel().GetRefWeapon()
}

func ModelGet(mt common.ModelType, id string) (*common.Model, error) {
	response := sendGetRequest(mt, id)
	err := response.GetErrorResponse()
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err.Error)
	}
	return response.GetGetResponse().Model, nil
}
