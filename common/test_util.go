package common

type TestUtil struct{}

func NewTestUtil() *TestUtil {
	tu := &TestUtil{}
	tu.Init()
	return tu
}

func (tu *TestUtil) Teardown() {
	original := Storage.(*ReadOnlyStorage).s
	Storage = original
}

func (tu *TestUtil) Init() {
	Storage = &ReadOnlyStorage{s: Storage}

	g1 := &Game{
		ID:   "my-cool-game",
		Name: "TEST",
		Players: []*Player{
			{Id: "john", Name: "John", GM: true},
			{Id: "par", Name: "Par Soulati"},
			{Id: "tim", Name: "tim"},
		},
	}
	ActiveGame = g1
	Games.Set(g1, g1.ID)
}

func (tu *TestUtil) BreastPlate() *Armor {
	a := &RefArmor{
		ID:         "breastplate",
		Name:       "Breastplate",
		CanDegrade: true,
		Material:   "Steel",
		SP:         8,
		Cost:       &MonetaryAmount{GP: 4},
		CanStack:   false,
		RefPenalty: -3,
	}

	References.Armors.Set(a, a.ID)

	armor := &Armor{
		ID:       GenerateID(),
		RefID:    a.ID,
		Equipped: true,
		SP:       a.SP,
		Degraded: false,
		Quality:  Quality_Quality_Standard,
	}
	return armor
}

func (tu *TestUtil) Sword() *Weapon {

	w := &RefWeapon{
		ID:                 "ref-weapon-1",
		Name:               "Sample Weapon",
		Description:        "Sample Description",
		Damage1H:           "2D6",
		Damage2H:           "3D6",
		BaseWeaponAccuracy: 1,
		DefenseModifier:    -3,
		DefenseType:        DefenseType_Parry,
		Cost:               &MonetaryAmount{SP: 75},
		Melee:              true,
		Rarity:             Rarity_Rarity_COMMON,
		RequiredSkill:      "bladed_weapons",
		Wield1Hand:         true,
		Wield2Hand:         true,
		CanBlock:           false,
		CanParry:           true,
		ParryMod:           -4,
		BlockMod:           0,
	}

	References.Weapons.Set(w, w.ID)

	weapon := &Weapon{
		ID:       GenerateID(),
		Carried:  true,
		Equipped: true,
		RefID:    w.ID,
		ParryMod: -999,
		BlockMod: -999,
	}
	return weapon
}
func (tu *TestUtil) BladedWeaponsSkill() *Skill {

	s := &RefSkill{
		ID:            "bladed_weapons",
		Name:          "Bladed Weapons",
		AttributeType: PrimaryAttributeVal_PrimaryAttributeVal_REF,
	}

	References.Skills.Set(s, s.ID)

	return &Skill{
		ID:    s.ID,
		Level: 5,
		IPs:   0,
	}
}

func (tu *TestUtil) Attributes(val int32) *CharacterAttributes {
	return &CharacterAttributes{
		BOD: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		REF: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		WILL: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		VIT: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		INT: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		LUCK: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		TECH: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		ATTR: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		ESS: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
		PER: &PrimaryAttribute{
			RawValue:  val,
			SetValue:  val,
			CalcValue: val,
		},
	}
}

func (tu *TestUtil) CharacterMelee(attrs int32) *Character {
	c1 := &Character{
		ID:         GenerateID(),
		Attributes: tu.Attributes(attrs),
		Weapons:    []*Weapon{tu.Sword()},
		Skills:     []*Skill{tu.BladedWeaponsSkill()},
		Armors:     []*Armor{tu.BreastPlate()},
	}

	ActiveGame.Characters().Set(c1, c1.ID)
	c1.Calculate()
	return c1
}
