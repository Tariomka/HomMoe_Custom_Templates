package registry

type receiversFilters struct {
	StartingHero string
	AllHeroes    string
}

// GetReceiversFilterValues returns the available receivers filter values used for
//
//	gameRules.bonuses.receiverFilter
func GetReceiversFilterValues() receiversFilters {
	return receiversFilters{
		StartingHero: "starting_hero",
		AllHeroes:    "all_heroes",
	}
}
