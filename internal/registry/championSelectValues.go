package registry

type championSelectRules struct {
	StartHero string
}

// GetChampionSelectValues returns the champion select rule types used for
//
//	gameRules.winConditions.championSelectRule
func GetChampionSelectValues() championSelectRules {
	return championSelectRules{
		StartHero: "StartHero",
	}
}
