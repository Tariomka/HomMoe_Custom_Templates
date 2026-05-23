package constants

import "strings"

type Victory struct {
	ID    string
	Label string
}

var VictoryCondition = struct {
	Standard         Victory
	LostStartingCity Victory
	HoldCity         Victory
	Tournament       Victory
}{
	Standard:         Victory{ID: "win_condition_1", Label: "Standard"},
	LostStartingCity: Victory{ID: "win_condition_3", Label: "Lost Starting City"},
	HoldCity:         Victory{ID: "win_condition_5", Label: "Hold City"},
	Tournament:       Victory{ID: "win_condition_6", Label: "Tournament"},
}

var VictoryConditions = []Victory{
	VictoryCondition.Standard,
	VictoryCondition.LostStartingCity,
	VictoryCondition.HoldCity,
	VictoryCondition.Tournament,
}

func GetVictoryCondition(id string) Victory {
	for _, victory := range VictoryConditions {
		if strings.EqualFold(victory.ID, id) {
			return victory
		}
	}
	return VictoryConditions[0]
}
