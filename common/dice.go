package common

import (
	"fmt"
	"strings"
)

// Properly checks and records
func Roll(dice *DiceRoll) *DiceRollResults {
	var result *DiceRollResults
	// Check if it needs to be asked for
	requestManual := false
	// If not then roll
	if requestManual {
		//
	} else {
		result = dice.Roll()
	}

	//TODO: Add to Game Log?

	return result
}

func ParseDiceRoll(roll string) *DiceRoll {
	return nil
}

func NewDieRoll(sides int32, exploding bool, reason DiceRollReason) *DiceRoll {
	d := &DiceRoll{
		Reason: reason,
		Dice: []*Die{
			{
				Sides:   sides,
				Amount:  1,
				Explode: exploding,
			},
		},
	}

	return d
}

func RollD10() *DiceRollResults {
	d := &DiceRoll{
		Dice: []*Die{
			{
				Sides:   10,
				Amount:  1,
				Explode: true,
			},
		},
	}

	return d.Roll()
}

func RollOnce(sides int32) int32 {
	return Random1(sides)
}

func (d *DiceRoll) Roll() *DiceRollResults {
	// Roll The Dice
	dieTotal := int32(0)
	results := make([]*DiceRollResult, len(d.Dice))
	for i, d := range d.Dice {
		rolls := make([]*DieRollResult, d.Amount)
		for n := 0; n < int(d.Amount); n++ {
			v := RollOnce(d.Sides)
			rolls[n] = &DieRollResult{
				Value: v,
			}
			dieTotal += v
		}
		results[i] = &DiceRollResult{
			Dice:    d,
			Results: rolls,
		}
	}

	// Sum the modifiers
	mods := int32(0)
	for _, m := range d.Modifiers {
		mods += m.Modifier
	}
	total := dieTotal + mods

	return &DiceRollResults{
		Rolls:     results,
		Total:     int32(total),
		Modifiers: d.Modifiers,
	}
}

func (d *Die) Format() string {
	return fmt.Sprintf("%vD%v", d.Amount, d.Sides)
}

func (d *DiceRoll) Format() string {
	// Roll The Dice
	var results []string
	for _, d := range d.Dice {
		results = append(results, d.Format())
	}

	// Sum the modifiers
	mods := int32(0)
	for _, m := range d.Modifiers {
		mods += m.Modifier
	}

	diceFmt := strings.Join(results, " + ")
	if mods > 0 {
		diceFmt = fmt.Sprintf("%v + %v", diceFmt, mods)
	} else if mods < 0 {
		absmods := -mods
		diceFmt = fmt.Sprintf("%v - %v", diceFmt, absmods)
	}
	return diceFmt
}

func (d *DiceRoll) AppendMod(tag string, mod int32) {
	d.Modifiers = append(d.Modifiers, &RollModifier{
		Modifier: mod,
		Tag:      tag,
	})
}

func (d *DiceRoll) FindModifier(tag string) int32 {
	for _, m := range d.Modifiers {
		if m.Tag == tag {
			return m.Modifier
		}
	}
	return 0
}

func (d *DiceRollResults) FindModifier(tag string) int32 {
	for _, m := range d.Modifiers {
		if m.Tag == tag {
			return m.Modifier
		}
	}
	return 0
}
