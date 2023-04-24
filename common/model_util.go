package common

func ToModels[T interface{}](items []T) []*Model {
	models := make([]*Model, len(items))
	for i, item := range items {
		models[i] = ToModel(item)
	}
	return models
}

// Holds Utiltity Methods for models
func ToModel(data interface{}) *Model {
	switch data.(type) {
	case *Armor:
		return &Model{
			Data: &Model_Armor{
				Armor: data.(*Armor),
			},
		}
	case *Character:
		return &Model{
			Data: &Model_Character{
				Character: data.(*Character),
			},
		}
	case *Game:
		return &Model{
			Data: &Model_Game{
				Game: data.(*Game),
			},
		}
	case *Player:
		return &Model{
			Data: &Model_Player{
				Player: data.(*Player),
			},
		}
	case *RefArmor:
		return &Model{
			Data: &Model_RefArmor{
				RefArmor: data.(*RefArmor),
			},
		}
	case *RefGameTerm:
		return &Model{
			Data: &Model_RefGameTerm{
				RefGameTerm: data.(*RefGameTerm),
			},
		}
	case *RefOrb:
		return &Model{
			Data: &Model_RefOrb{
				RefOrb: data.(*RefOrb),
			},
		}
	case *RefSkill:
		return &Model{
			Data: &Model_RefSkill{
				RefSkill: data.(*RefSkill),
			},
		}
	case *RefWeapon:
		return &Model{
			Data: &Model_RefWeapon{
				RefWeapon: data.(*RefWeapon),
			},
		}
	}
	return nil
}
