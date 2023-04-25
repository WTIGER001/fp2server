package common

import (
	"log"
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
	encounter.DebugPrintInitiativeOrder()

	// Round 1..
	round := encounter.NextRound()
	encounter.DebugPrintRound(round)

	turnIndex := -1
	for encounter.NextTurn() {
		// Print
		log.Printf("New Turn: %v", encounter.CurrentTurn)
		turn := encounter.GetTurn()
		if turn.Status == TurnStatus_TurnStatus_Pending {
			turn.Status = TurnStatus_TurnStatus_Held
		}

		switch turn.CharacterId {
		case c1.ID:
			log.Printf("%v - %v - %v", turn.Order, turn.Status, c1.Name)
		case c2.ID:
			log.Printf("%v - %v - %v", turn.Order, turn.Status, c2.Name)
		}

		turnIndex++
		assert.Equal(t, encounter.CurrentRound, int32(0))
		assert.GreaterOrEqual(t, encounter.CurrentTurn, int32(turnIndex))
	}
	log.Printf("END OF ROUND")
	encounter.DebugPrintRound(round)

}
