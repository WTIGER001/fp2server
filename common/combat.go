package common

// Handles Processing for Combat
type CombatProcessor struct {
}

// Handles the messages relevant to combat.
func (cp *CombatProcessor) Handle(m *Fp2Message) {
	switch m.Data.(type) {

	case *Fp2Message_Attack:
		cp.OnAttack(m.GetAttack())

	case *Fp2Message_AttackResult:

	case *Fp2Message_DefenseChallenge:

	case *Fp2Message_DefenseChallengeResponse:

	}

}

func (cp *CombatProcessor) OnAttack(attack *Attack) {
	// Get Defense Option
	cp.GetDefense(attack)
}

/*
	TODO: Support Stacking Armor
	FIXME: Look up the Defensive Item
	TODO: Allow for strange magical effects that actively defend
	FIXME: ROLL Correctly
*/
func (cp *CombatProcessor) ResolveAttack(attack *Attack, defense *DefenseOption) *AttackResult {
	// Roll Attack
	attackRoll := RollOnce(10, true)

	// Roll Defense
	defenseRoll := RollOnce(10, true)
	defenseTotal := defenseRoll + defense.DefenseTotal
	dr := AttackMeleeDR
	if defenseTotal < dr {
		defenseTotal = dr
	}

	result := &AttackResult{
		Attack:       attack,
		Defense:      defense,
		AttackTotal:  attackRoll,
		DefenseTotal: defenseTotal,
		ArmorSP:      10,
	}

	// Resolve
	if attackRoll >= defenseTotal {
		damage := RollOnce(10, true)

		target, _ := ActiveGame.Characters().Get(attack.Target)
		armor := target.GetArmor()

		// Success
		inflicted := damage - armor.SP
		if inflicted > 0 {
			// Check if decrement
			if armor.CanDegradeDefault() {
				result.ArmorDamaged = true
			}
			result.DamageInflicted = inflicted
		}
	} else if defense.DefenseType == DefenseType_Block {
		damage := RollOnce(10, true)

		target, _ := ActiveGame.Characters().Get(attack.Target)
		armor := target.GetArmor()
		shield := target.GetArmor()

		inflicted := damage - shield.SP
		if inflicted > 0 {
			// Check if decrement
			if shield.CanDegradeDefault() {
				result.ShieldDamaged = true
			}

			inflicted2 := inflicted - armor.SP
			if inflicted2 > 0 {
				// Check if decrement
				if armor.CanDegradeDefault() {
					result.ArmorDamaged = true
				}
				result.DamageInflicted = inflicted
			}
		}
	} else if defense.DefenseType == DefenseType_Parry {

	} else if defense.DefenseType == DefenseType_Dodge {

	}

	return result
}

func (cp *CombatProcessor) GetDefense(attack *Attack) *DefenseOption {
	// Construct Defense Options
	options := cp.GetDefenseOptions(attack)

	// Get the Defensive Choice
	choice := cp.GetDefenseChoice(attack, options)

	return choice
}

func (cp *CombatProcessor) GetDefenseChoice(attack *Attack, options []*DefenseOption) *DefenseOption {
	// This can calculate the choice
	return options[0]
}

// The Defensive options for a character are based on their weapons, skills, sheilds, etc.
func (cp *CombatProcessor) GetDefenseOptions(attack *Attack) []*DefenseOption {
	target, _ := ActiveGame.Characters().Get(attack.Target)
	var opts []*DefenseOption

	// Handles Weapons and Shields
	if target.Weapons != nil {
		for _, w := range target.Weapons {
			// Get the reference
			wRef, _ := References.Weapons.Get(w.RefID)

			// Get the Skill
			skillTotal := target.GetCurrentSkillTotal(wRef.RequiredSkill)

			// Handle Missing
			if wRef.CanParry {
				// Calculate the Name
				name := First(w.Name, wRef.Name)

				//Calculate the Parry
				parryMod := FirstInt32(w.ParryMod, wRef.ParryMod)

				opt := &DefenseOption{
					DefenseType:  DefenseType_Parry,
					DefenseItem:  w.ID,
					Description:  Text.Format(TXT_COMBAT_PARRY, name),
					DefenseTotal: skillTotal + parryMod,
				}
				opts = append(opts, opt)
			}

			if wRef.CanBlock {
				// Calculate the Name
				name := First(w.Name, wRef.Name)

				//Calculate the Block
				blockMod := FirstInt32(w.BlockMod, wRef.BlockMod)

				opt := &DefenseOption{
					DefenseType:  DefenseType_Block,
					DefenseItem:  w.ID,
					Description:  Text.Format(TXT_COMBAT_BLOCK, name),
					DefenseTotal: skillTotal + blockMod,
				}
				opts = append(opts, opt)
			}
		}
	}

	// Always Add a Do Nothing
	doNothing := &DefenseOption{
		DefenseType: DefenseType_None,
		Description: "Accept the Attack",
	}
	opts = append(opts, doNothing)
	return opts
}

// Rolls the Initiave for a Character
// func RollInititative(c *Character) int32 {
// 	base := c.Attributes.Initiative.Value
// 	r := &DiceRoll{
// 		Dice: []*Die{
// 			{
// 				Amount:  1,
// 				Sides:   10,
// 				Explode: true,
// 			},
// 		},
// 		Modifiers: []*RollModifier{
// 			{
// 				Modifier: base,
// 			},
// 		},
// 	}

// 	result := RequestRoll("encounter.Initiative", r, c)
// 	return result.Total
// }
