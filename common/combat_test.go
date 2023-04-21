package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttack1(t *testing.T) {
	tu := &TestUtil{}
	Storage = &ReadOnlyStorage{s: Storage}

	c1 := &Character{
		ID:         GenerateID(),
		Attributes: tu.Attributes(7),
		Weapons:    []*Weapon{tu.Sword()},
		Skills:     []*Skill{tu.BladedWeaponsSkill()},
		Armors:     []*Armor{tu.BreastPlate()},
	}
	c2 := &Character{
		ID:         GenerateID(),
		Attributes: tu.Attributes(10),
		Weapons:    []*Weapon{tu.Sword()},
		Skills:     []*Skill{tu.BladedWeaponsSkill()},
		Armors:     []*Armor{tu.BreastPlate()},
	}
	References.Characters.Set(c1, c1.ID)
	References.Characters.Set(c2, c2.ID)

	cp := &CombatProcessor{}
	attack := &Attack{
		Attacker:    c1.ID,
		Target:      c2.ID,
		AttackType:  AttackType_AttackType_Melee,
		Description: "",
		Weapon:      c2.Weapons[0].ID,
	}

	choice := cp.GetDefense(attack)
	assert.Equal(t, choice.DefenseType, DefenseType_Parry)

}
