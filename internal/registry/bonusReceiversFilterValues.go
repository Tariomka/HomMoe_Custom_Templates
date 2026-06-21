package registry

type receiversFilters struct {
	StartingHero string
	AllHeroes    string
}

var bonusReceiversFilterValues = receiversFilters{
	StartingHero: "starting_hero",
	AllHeroes:    "all_heroes",
}

// GetReceiversFilterValues returns the available receivers filter values used for
//
//	gameRules.bonuses.receiverFilter
func GetReceiversFilterValues() receiversFilters {
	return bonusReceiversFilterValues
}
