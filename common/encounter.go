package common

import (
	"sort"
)

// Roll the initative for an encounter. This should only be done once
// per encounter. If you need to add or remove a person from the
// encounter there are other methods.
func (e *Encounter) RollInititative() {
	// Roll the initatives
	var orders []*InitiativeOrder
	for i, entityRef := range e.Entities {
		// Look up the correct entity
		entity := e.GetEntity(e.ID)
		// Get the raw value
		raw := entity.Attributes.Initiative.Value
		// Roll the initative dice
		dice := Clone(InitiativeDice)
		dice.AppendMod("Initiative", raw)
		result := Roll(dice)
		// Add to temp
		orders = append(orders, &InitiativeOrder{
			EntityID:        entityRef.ID,
			Value:           result.Total,
			DiceRollResults: result,
		})
	}

	// Sort the initative rolls
	sort.Slice(orders, func(i, j int) bool {
		o1 := orders[i]
		o2 := orders[j]
		return InitativeOrderLess(o1, o2)
	})

	// Assign
	e.InitiativeOrders = orders

	// initialize
	e.CurrentRound = -1
	e.CurrentTurn = -1
	e.Rounds = nil
}

func InitativeOrderLess(o1 *InitiativeOrder, o2 *InitiativeOrder) bool {
	// Level 1
	if o1.Value < o2.Value {
		return true
	}
	if o1.Value > o2.Value {
		return false
	}

	// Level 2, base values
	b1 := o1.DiceRollResults.FindModifier("Initiative")
	b2 := o2.DiceRollResults.FindModifier("Initiative")
	if b1 < b2 {
		return true
	}
	if b1 > b2 {
		return false
	}

	// Level3, Luck?
	// TBD : Add luck here

	// Level 4: Random
	for {
		r1 := Random64()
		r2 := Random64()
		if r1 < r2 {
			return true
		}
		if r1 > r2 {
			return false
		}
	}
}

func (e *Encounter) GetEntity(id string) *Character {
	entity, _ := ActiveGame.Characters().Get(o.EntityID)
	return entity
}

// BuildNextRound builds out the action order
// for the next round and places it in the array
// It tries to read the previous round to see
// if there are any borrowed reactions.
func (e *Encounter) BuildNextRound() {
	// Start at the beginning of the initative order
	// Just loop through the initative order and construct
	// the order object.
	round := &Round{}
	round.RoundNumber = int32(len(e.Rounds) + 1)

	for i, o := range e.InitiativeOrders {
		// Get the Entity that is going. This can be a Character
		// Creature, or NPC. If this is a character then it should
		// also be in the ActiveGame object. If it is a creature or
		// NPC then it is either in the Game Object or in the references
		// FIXME: Support creatures and NPCs
		entity := e.GetEntity(o.EntityID)
		raw := o.Value
		action := int32(0)
		for ; action < entity.ActionCount; action++ {
			turn := &Turn{
				// Assign the same id
				CharacterId: o.EntityID,

				// Calculate the order. This is the raw initiative, decremented
				// by SupsequentInitativeActionPenalty (3) and then, so things sort
				// correctly add the index order / 10 (as a decimal palce)
				Order: float32(raw-(action*SupsequentInitativeActionPenalty)) + float32(i/10.0),

				// Set the default action to just be a single
				Actions: 1,

				// Initally start off with all actions as PENDING
				Status: TurnStatus_TurnStatus_Pending,
			}
			round.Turns = append(round.Turns, turn)
		}
	}

	// The inital order is calculated. Now check for any borrowed actions in
	// the preivous round
	if len(e.Rounds) > 0 {
		// Get the last round
		previous := e.Rounds[len(e.Rounds)-1]
		for _, reaction := range previous.DefensiveReactions {
			turn := round.FindNextPending(reaction.EntityID)
			if turn != nil {
				turn.Status = TurnStatus_TurnStatus_Borrowed
			}
		}
	}

	// Now add the round
	e.Rounds = append(e.Rounds, round)
}

func (e *Encounter) ExpireEffects() {

}

// Records a reaction from a character in the
// round. This will first draw from all the
// available defensive reactions and then will
// draw from the held or pending actions remaining
// if the character does not have anything to draw
// from then "false" is returned
func (e *Encounter) CharacterReact(id string) bool {
	// count the defensive reactions used in this turn
	// if there are any remaining then use them
	entity := e.GetEntity(id)
	round := e.GetCurrentRound()
	defensiveReactionsUsed := 0
	for _, dr := range round.DefensiveReactions {
		if dr.EntityID == id {
			defensiveReactionsUsed++
		}
	}
	if defensiveReactionsUsed < int(entity.DefensiveReactions) {
		round.DefensiveReactions = append(round.DefensiveReactions, &DefensiveReaction{
			EntityID:  id,
			TurnIndex: e.CurrentTurn,
		})
		return true
	}

	// No Defensive Reactions left, instead use a pending one
	pendingTurn := round.FindNextPending(id)
	if pendingTurn != nil {
		pendingTurn.Status = TurnStatus_TurnStatus_Reacted
	}

	// No Actions this round.. Try to borrow from the next round
	borrowedReactionsUsed := 0
	for _, dr := range round.BorrowedReactions {
		if dr.EntityID == id {
			borrowedReactionsUsed++
		}
	}
	if defensiveReactionsUsed < AllowBorrowedReactions {
		round.BorrowedReactions = append(round.BorrowedReactions, &DefensiveReaction{
			EntityID:  id,
			TurnIndex: e.CurrentTurn,
		})
		return true
	}

	// No MORE
	return false
}

func (e *Encounter) GetCurrentRound() *Round {
	if len(e.Rounds) == 0 {
		return nil
	}
	return e.Rounds[len(e.Rounds)-1]
}

// Advances the round to the next round. This will generate
// the next turn order, clear any status effects and go to the
// first real turn
func (e *Encounter) NextRound() {
	if len(e.InitiativeOrders) == 0 {
		// MAYBE Raise an error here?
		return
	}
	e.BuildNextRound()
	e.ExpireEffects()
	e.CurrentRound = int32(len(e.Rounds) - 1)
	e.CurrentTurn = -1 // Not Started
	e.NextTurn()
}

// Advances to the next active turn
func (e *Encounter) NextTurn() {
	if e.CurrentTurn == 999 {
		// Cannot advance until a new round
		return
	}
	// Stupid checks
	if len(e.InitiativeOrders) == 0 {
		return
	}
	if len(e.Rounds) == 0 {
		e.BuildNextRound()
	}

	// Get the current Round
	round := e.Rounds[len(e.Rounds)-1]

	// See if this was the LAST Turn
	nextTurnIndex := e.CurrentTurn + 1
	if len(round.Turns) >= int(nextTurnIndex) {
		// TIme to make a new turn
		e.CurrentTurn = 99
		// Wait for someone to manually advance the round
		return
	}
	e.CurrentTurn = nextTurnIndex

	// Now last minute apply, because the previous action could
	// have changed things.
	e.ApplyStatusConditons()

	// Now, check the upcoming turn to see if it is a "real turn"
	turn := round.Turns[e.CurrentTurn]
	switch turn.Status {
	case TurnStatus_TurnStatus_Acted:
		// UMM Error, this should only occur AFTER
	case TurnStatus_TurnStatus_Borrowed:
		// Skip, dont care about this
		e.NextTurn()
	case TurnStatus_TurnStatus_Held:
		// UMM Error, this should only occur AFTER
	case TurnStatus_TurnStatus_Pending:
		// On a GOOD State
		return
	case TurnStatus_TurnStatus_Reacted:
		// Skip, dont care about this
		e.NextTurn()
	case TurnStatus_TurnStatus_Reacted_Borrowed:
		// Skip, dont care about this
		e.NextTurn()
	}
}

// FindNextPending finds the next pending turn for an entity
func (r *Round) FindNextPending(entityId string) *Turn {
	for _, t := range r.Turns {
		if t.CharacterId == entityId && t.Status == TurnStatus_TurnStatus_Pending {
			return t
		}
	}
	return nil
}

func (e *Encounter) AddEntities(newpeople []*EntityReference) {

}

func (e *Encounter) RemoveEntities(toRemove []*EntityReference) {

}

func (e *Encounter) ApplyStatusConditons() {

}