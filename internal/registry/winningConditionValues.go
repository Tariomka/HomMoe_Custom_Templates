package registry

type winConditions struct {
	Standard    string
	CapitalHold string
	FinalBattle string
	CityHold    string
	Tournament  string
}

var winningConditionValues = winConditions{
	Standard: "win_condition_1",
	// "win_condition_2"
	CapitalHold: "win_condition_3",
	FinalBattle: "win_condition_4",
	CityHold:    "win_condition_5",
	Tournament:  "win_condition_6",
}

// GetWinningConditionValues returns the available winning condition types used for
//
//	displayWinCondition
func GetWinningConditionValues() winConditions {
	return winningConditionValues
}
