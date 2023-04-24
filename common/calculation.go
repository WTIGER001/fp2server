package common

type Calculated interface {
	Calculate()
}

// Calculates all the fields for a character. This should be done sparingly.
// There are optimizations that can be done later.
func (c *Character) Calculate() {

	// Need to determine the correct order.

	// Calculate all the modifications for all active effects. This is weapons equipped,
	//    Armor equipped, Gear equipped, Targeted Ablitilies, Potions (etc.)
	mods := c.AllActiveModSums2()

	// Calculate "Permenant" Attributes
	c.Attributes.ATTR.SetValue = c.Attributes.ATTR.RawValue + mods[ModificationType_MT_AttributeATTR][""]
	c.Attributes.BOD.SetValue = c.Attributes.BOD.RawValue + mods[ModificationType_MT_AttributeBOD][""]
	c.Attributes.ESS.SetValue = c.Attributes.ESS.RawValue + mods[ModificationType_MT_AttributeESS][""]
	c.Attributes.INT.SetValue = c.Attributes.INT.RawValue + mods[ModificationType_MT_AttributeINT][""]
	c.Attributes.REF.SetValue = c.Attributes.REF.RawValue + mods[ModificationType_MT_AttributeREF][""]
	c.Attributes.LUCK.SetValue = c.Attributes.LUCK.RawValue + mods[ModificationType_MT_AttributeLUCK][""]
	c.Attributes.WILL.SetValue = c.Attributes.WILL.RawValue + mods[ModificationType_MT_AttributeWILL][""]
	c.Attributes.TECH.SetValue = c.Attributes.TECH.RawValue + mods[ModificationType_MT_AttributeTECH][""]
	c.Attributes.PER.SetValue = c.Attributes.PER.RawValue + mods[ModificationType_MT_AttributePER][""]
	c.Attributes.VIT.SetValue = c.Attributes.VIT.RawValue + mods[ModificationType_MT_AttributeVIT][""]

	// Calculate "Temp" Attributes
	c.Attributes.ATTR.CalcValue = c.Attributes.ATTR.SetValue + mods[ModificationType_MT_AttributeTempATTR][""]
	c.Attributes.BOD.CalcValue = c.Attributes.BOD.SetValue + mods[ModificationType_MT_AttributeTempBOD][""]
	c.Attributes.ESS.CalcValue = c.Attributes.ESS.SetValue + mods[ModificationType_MT_AttributeTempESS][""]
	c.Attributes.INT.CalcValue = c.Attributes.INT.SetValue + mods[ModificationType_MT_AttributeTempINT][""]
	c.Attributes.REF.CalcValue = c.Attributes.REF.SetValue + mods[ModificationType_MT_AttributeTempREF][""]
	c.Attributes.LUCK.CalcValue = c.Attributes.LUCK.SetValue + mods[ModificationType_MT_AttributeTempLUCK][""]
	c.Attributes.WILL.CalcValue = c.Attributes.WILL.SetValue + mods[ModificationType_MT_AttributeTempWILL][""]
	c.Attributes.TECH.CalcValue = c.Attributes.TECH.SetValue + mods[ModificationType_MT_AttributeTempTECH][""]
	c.Attributes.PER.CalcValue = c.Attributes.PER.SetValue + mods[ModificationType_MT_AttributeTempPER][""]
	c.Attributes.VIT.CalcValue = c.Attributes.VIT.SetValue + mods[ModificationType_MT_AttributeTempVIT][""]

	// Apply wound levels
	wounds := c.GetWoundState()
	if wounds == WoundState_WoundState_Serious {
		c.Attributes.REF.CalcValue = c.Attributes.REF.CalcValue - 2
	} else if wounds == WoundState_WoundState_Critical {
		c.Attributes.REF.CalcValue = c.Attributes.REF.CalcValue / 2
		c.Attributes.WILL.CalcValue = c.Attributes.WILL.CalcValue / 2
		c.Attributes.INT.CalcValue = c.Attributes.INT.CalcValue / 2
	} else if wounds == WoundState_WoundState_Mortal {
		c.Attributes.REF.CalcValue = c.Attributes.REF.CalcValue / 3
		c.Attributes.WILL.CalcValue = c.Attributes.WILL.CalcValue / 3
		c.Attributes.INT.CalcValue = c.Attributes.INT.CalcValue / 3
	}

	// Skills
	for _, s := range c.Skills {
		skillRef, _ := References.Skills.Get(s.ID)
		attr := skillRef.AttributeType
		attrVal := c.GetCurentAttribute(attr)

		s.Mod = mods[ModificationType_MT_Skill][s.ID]
		s.Total = s.Level + s.Mod + attrVal
	}

	// Calculate Derived
	if c.Attributes.Initiative == nil {
		c.Attributes.Initiative = &CalculatedValue{}
	}
	if c.Attributes.Awarness == nil {
		c.Attributes.Awarness = &CalculatedValue{}
	}
	c.Attributes.Initiative.Value = c.Attributes.PER.CalcValue + mods[ModificationType_MT_Initiative][""]
	c.Attributes.Awarness.Value = c.Attributes.PER.CalcValue + mods[ModificationType_MT_Awarness][""]
	c.ActionCount = DefaultActions + mods[ModificationType_MT_Actions][""]
	c.DefensiveReactions = DefaultDefensiveReactions + mods[ModificationType_MT_DefensiveReactions][""]
}
