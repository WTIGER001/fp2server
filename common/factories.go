package common

import "google.golang.org/protobuf/proto"

type ModelFactory[T proto.Message] interface {
	New() T
}

var CharacterFactory = func() *Character { return &Character{} }
var EncounterFactory = func() *Encounter { return &Encounter{} }
var GameFactory = func() *Game { return &Game{} }
var WeaponFactory = func() *Weapon { return &Weapon{} }
var ArmorFactory = func() *Armor { return &Armor{} }
var OrbFactory = func() *Orb { return &Orb{} }

var RefWeaponFactory = func() *RefWeapon { return &RefWeapon{} }
var RefGameTermFactory = func() *RefGameTerm { return &RefGameTerm{} }
var RefSkillFactory = func() *RefSkill { return &RefSkill{} }
var RefArmorFactory = func() *RefArmor { return &RefArmor{} }
var RefOrbFactory = func() *RefOrb { return &RefOrb{} }

// var RefGearFactory = func() *RefGear { return &RefGear{} }
