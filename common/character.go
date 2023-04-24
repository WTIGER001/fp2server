package common

// Rolls 1d10 and then returns the value. If the roll
// is above a certain number then increase the IPs.
func (c *Character) UseSkill(ID string) int32 {

	roll := RollD10()
	d10 := roll.Rolls[0].Results[0].Value
	refSkill, _ := References.Skills.Get(ID)
	if refSkill == nil {
		return -1
	}

	s := c.FindSkill(ID)
	if d10 >= IPGainThreshold {
		s.IPs += 1
	}
	return s.Total + d10
}

func (c *Character) GetWoundState() WoundState {
	if c.Health.CurrentWounds == 0 {
		return WoundState_WoundState_None
	}
	if c.Health.CurrentWounds <= c.Health.LightLevels {
		return WoundState_WoundState_Light
	}
	if c.Health.CurrentWounds <= c.Health.LightLevels+c.Health.SeriousLevels {
		return WoundState_WoundState_Serious
	}
	if c.Health.CurrentWounds <= c.Health.LightLevels+c.Health.SeriousLevels+c.Health.CriticalLevels {
		return WoundState_WoundState_Critical
	}
	return WoundState_WoundState_Mortal
}

func (c *Character) GetCurrentSkillTotal(ID string) int32 {
	// Look for the Skill
	refSkill, _ := References.Skills.Get(ID)
	if refSkill == nil {
		return 0
	}

	s := c.FindSkill(ID)
	// attr := c.GetCurentAttribute(refSkill.AttributeType)
	// // Skill Total = Level + Attribute + MODs
	// total := s.Level + attr + s.Mod
	// return total

	// Now what to do about skills with a level of 0
	if s == nil {
		return 0
	} else {
		return s.Total
	}
}

func (c *Character) FindSkill(ID string) *Skill {
	for _, s := range c.Skills {
		if s.ID == ID {
			return s
		}
	}
	return nil
}

func (c *Character) GetCurentAttribute(attr PrimaryAttributeVal) int32 {
	switch attr {
	case PrimaryAttributeVal_PrimaryAttributeVal_ATTR:
		return c.Attributes.ATTR.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_BOD:
		return c.Attributes.BOD.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_ESS:
		return c.Attributes.ESS.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_INT:
		return c.Attributes.INT.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_LUCK:
		return c.Attributes.LUCK.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_PER:
		return c.Attributes.PER.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_REF:
		return c.Attributes.REF.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_TECH:
		return c.Attributes.TECH.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_WILL:
		return c.Attributes.WILL.CalcValue
	case PrimaryAttributeVal_PrimaryAttributeVal_UNK:
	}
	return 0
}

func (c *Character) GetArmor() *Armor {
	for _, a := range c.Armors {
		if a.Equipped {
			return a
		}
	}

	return nil
}

func (c *Character) InitHealth() {
	if c.Health == nil {
		c.Health = &CharacterHealth{}
	}

	// Calculate Base Levels
	bod := c.Attributes.BOD.RawValue
	c.Health.LightLevels = HealthLevelsFromBod(bod) + c.SumActiveMods(ModificationType_MT_WoundLevelsLight)
	c.Health.SeriousLevels = HealthLevelsFromBod(bod) + c.SumActiveMods(ModificationType_MT_WoundLevelsSerious)
	c.Health.CriticalLevels = HealthLevelsFromBod(bod) + c.SumActiveMods(ModificationType_MT_WoundLevelsCritical)
	c.Health.MortalLevels = HealthLevelsFromBod(bod) + c.SumActiveMods(ModificationType_MT_WoundLevelsMortal)

	// Look for ANY items that provide ability to increase / decrease levels
	// TODO: Need to standardize this... Likely need a "Type" and a "Mod"
	//      This would apply to wound levels, attributes, skills, weapons,
	//		armor, etc.

}

func (c *Character) InflictDamage(amount int32) {
	c.InitHealth()

	// Apply the damage
	c.Health.CurrentWounds = c.Health.CurrentWounds + amount

}

func (c *Character) HealDamage() {

}

func (c *Character) GetAllActiveMods(t ModificationType) []*Modification {
	var rtn []*Modification

	// Weapons
	for _, w := range c.Weapons {
		if w.Equipped {
			for _, m := range w.Modifications {
				if m.Type == t {
					rtn = append(rtn, m)
				}
			}
		}
	}

	// Armor
	for _, a := range c.Armors {
		if a.Equipped {
			for _, m := range a.Modifications {
				if m.Type == t {
					rtn = append(rtn, m)
				}
			}
		}
	}

	return rtn
}

func (c *Character) SumActiveMods(t ModificationType) int32 {
	var total int32
	active := c.GetAllActiveMods(t)
	for _, m := range active {
		total += m.Amount
	}

	return total
}

func (c *Character) AllActiveModSums() map[ModificationType]int32 {
	rtn := make(map[ModificationType]int32)

	// Weapons
	for _, w := range c.Weapons {
		if w.Equipped {
			for _, m := range w.Modifications {
				rtn[m.Type] = rtn[m.Type] + m.Amount
			}
		}
	}

	// Armor
	for _, a := range c.Armors {
		if a.Equipped {
			for _, m := range a.Modifications {
				rtn[m.Type] = rtn[m.Type] + m.Amount
			}
		}
	}

	return rtn
}

func (c *Character) AllActiveModSums2() map[ModificationType]map[string]int32 {
	rtn := make(map[ModificationType]map[string]int32)

	// Weapons
	for _, w := range c.Weapons {
		if w.Equipped {
			for _, m := range w.Modifications {
				cats := rtn[m.Type]
				if cats == nil {
					cats = make(map[string]int32)
					rtn[m.Type] = cats
				}
				cats[m.IDAffected] = cats[m.IDAffected] + m.Amount
			}
		}
	}

	// Armor
	for _, a := range c.Armors {
		if a.Equipped {
			for _, m := range a.Modifications {
				cats := rtn[m.Type]
				if cats == nil {
					cats = make(map[string]int32)
					rtn[m.Type] = cats
				}
				cats[m.IDAffected] = cats[m.IDAffected] + m.Amount
			}
		}
	}

	return rtn
}
