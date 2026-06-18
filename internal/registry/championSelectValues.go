package registry

type championSelectRules struct {
	StartHero string
}

var championSelectValues = championSelectRules{
	StartHero: "StartHero",
}

// GetChampionSelectValues returns the champion select rule types used for
//
//	gameRules.winConditions.championSelectRule
func GetChampionSelectValues() championSelectRules {
	return championSelectValues
}
