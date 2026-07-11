package registry

type bonusesSids struct {
	HeroAdditionalUnits string
	HeroExperience      string
	HeroItem            string
	HeroSpell           string
	HeroStat            string
	HeroUnitMultiplier  string
	Resource            string
	SideExperience      string
}

// GetMapBonusesValues returns the available map bonus SIDs used for
//
//	gameRules.bonuses.sid
func GetMapBonusesValues() bonusesSids {
	return bonusesSids{
		HeroItem:            "add_bonus_hero_item",
		Resource:            "add_bonus_res",
		HeroExperience:      "add_bonus_hero_exp",
		SideExperience:      "add_bonus_side_exp",
		HeroSpell:           "add_bonus_hero_spell",
		HeroUnitMultiplier:  "add_bonus_hero_unit_multipler",
		HeroAdditionalUnits: "add_bonus_hero_unit",
		HeroStat:            "add_bonus_hero_stat",
	}
}
