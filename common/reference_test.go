package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestProtoMarshall(t *testing.T) {
	w := &RefWeapon{
		ID:                 "ref-weapon-1",
		Name:               "Sample Weapon",
		Description:        "Sample Description",
		Damage1H:           "2D6",
		Damage2H:           "3D6",
		BaseWeaponAccuracy: 1,
		DefenseModifier:    -3,
		DefenseType:        DefenseType_Parry,
		Cost:               &MonetaryAmount{SP: 75},
		Melee:              true,
		Rarity:             Rarity_Rarity_COMMON,
		RequiredSkill:      "bladed_weapons",
		Wield1Hand:         true,
		Wield2Hand:         true,
	}
	data, err := proto.Marshal(w)
	assert.Nil(t, err)

	w2 := &RefWeapon{}
	err = proto.Unmarshal(data, w2)

	assert.Nil(t, err)
	assert.Equal(t, w.ID, w2.ID)

	err = os.WriteFile("e:\\fp2\\test.pbf", data, 0600)
	assert.Nil(t, err)

	data2, err := os.ReadFile("e:\\fp2\\test.pbf")
	assert.Nil(t, err)

	assert.ElementsMatch(t, data, data2)

	w3 := &RefWeapon{}
	err = proto.Unmarshal(data, w3)
	assert.Nil(t, err)

	data3, err := os.ReadFile("e:\\fp2\\ref-weapon-1.pbf")
	assert.Nil(t, err)

	assert.ElementsMatch(t, data, data3)

	w4 := &RefWeapon{}
	err = proto.Unmarshal(data, w4)
	assert.Nil(t, err)

	References.Weapons.Set(w4, w4.ID)
	w5 := References.Weapons.Load(w4.ID)
	assert.NotNil(t, w5)

}

func TestReference(t *testing.T) {
	References.LoadAll()

	cntWeapons := References.Weapons.Count()
	assert.Equal(t, 0, cntWeapons)

	w := &RefWeapon{
		ID:                 "ref-weapon-1",
		Name:               "Sample Weapon",
		Description:        "Sample Description",
		Damage1H:           "2D6",
		Damage2H:           "3D6",
		BaseWeaponAccuracy: 1,
		DefenseModifier:    -3,
		DefenseType:        DefenseType_Parry,
		Cost:               &MonetaryAmount{SP: 75},
		Melee:              true,
		Rarity:             Rarity_Rarity_COMMON,
		RequiredSkill:      "bladed_weapons",
		Wield1Hand:         true,
		Wield2Hand:         true,
	}

	References.Weapons.Set(w, w.ID)

}

func TestReferenceGet(t *testing.T) {
	w := References.Weapons.Get("ref-weapon-1")
	assert.NotNil(t, w)
}
