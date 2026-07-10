package registry

type spellSchoolTypes struct {
	HighNeutral string
	Daylight    string
	Nightshade  string
	Arcane      string
	Primal      string
}

func GetSpellSchoolTypeValues() spellSchoolTypes {
	return spellSchoolTypes{
		HighNeutral: "neutral",
		Daylight:    "day",
		Nightshade:  "night",
		Arcane:      "space",
		Primal:      "primal",
	}
}

func GetSpellSchoolTypeList() []string {
	spellSchoolValues := GetSpellSchoolTypeValues()
	return []string{
		spellSchoolValues.HighNeutral,
		spellSchoolValues.Daylight,
		spellSchoolValues.Nightshade,
		spellSchoolValues.Arcane,
		spellSchoolValues.Primal,
	}
}
