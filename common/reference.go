package common

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	"google.golang.org/protobuf/proto"
)

var References *ReferenceManager

type ReferenceItemManager[T proto.Message] struct {
	manager      *ReferenceManager
	refType      ReferenceType
	path         string
	cache        cmap.ConcurrentMap[string, T]
	loadingAll   *sync.WaitGroup
	newItemFn    func() T
	unmarshallFn func(data []byte) (T, error)
}

type ReferenceClient struct {
}

type ReferenceManager struct {
	loadingAll *sync.WaitGroup

	Weapons   *ItemManager[*RefWeapon]
	GameTerms *ItemManager[*RefGameTerm]
	Skills    *ItemManager[*RefSkill]
	Armors    *ItemManager[*RefArmor]
	Orbs      *ItemManager[*RefOrb]
}

func (r *ReferenceManager) LoadAll() {
	r.Weapons.GetAll()
	r.Skills.GetAll()
	r.GameTerms.GetAll()
	r.Armors.GetAll()
	r.Orbs.GetAll()
}

func newReferenceManager() *ReferenceManager {
	r := &ReferenceManager{}
	r.loadingAll = new(sync.WaitGroup)

	r.Weapons = NewItemManager("references/weapons", RefWeaponFactory)
	r.GameTerms = NewItemManager("references/gameterms", RefGameTermFactory)
	r.Skills = NewItemManager("references/skills", RefSkillFactory)
	r.Armors = NewItemManager("references/armors", RefArmorFactory)
	r.Orbs = NewItemManager("references/orbs", RefOrbFactory)

	return r
}
