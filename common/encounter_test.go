package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncounter1(t *testing.T) {
	tu := NewTestUtil()
	defer tu.Teardown()

	c1 := tu.CharacterMelee(7)
	c2 := tu.CharacterMelee(10)

	characters := []string{c1.ID, c2.ID}
	encounter := NewEncounter(characters, nil)
	encounter.Name = "TestEncounter1"
	ActiveGame.Encounters().Set(encounter, encounter.ID)

	encounter.RollInititative()
	assert.Equal(t, len(encounter.InitiativeOrders), 2)

	// Round 1..
	encounter.NextRound()

	turn := -1
	for encounter.NextTurn() {
		turn++
		assert.Equal(t, encounter.CurrentRound, 0)
		assert.GreaterOrEqual(t, encounter.CurrentTurn, turn)
	}

}
