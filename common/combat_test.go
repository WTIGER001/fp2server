package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttack1(t *testing.T) {
	tu := NewTestUtil()
	defer tu.Teardown()

	c1 := tu.CharacterMelee(7)
	c2 := tu.CharacterMelee(10)

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
