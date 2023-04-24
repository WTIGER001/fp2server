package common

func (a *Armor) CanDegradeDefault() bool {
	if a.CanDegrade == BooleanValue_BooleanValue_True {
		return true
	}

	if a.CanDegrade == BooleanValue_BooleanValue_False {
		return false
	}

	armorRef, _ := References.Armors.Get(a.RefID)
	return armorRef.CanDegrade
}
