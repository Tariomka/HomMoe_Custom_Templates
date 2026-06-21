package constants

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var winConditions = registry.GetWinningConditionValues()

type Victory struct {
	ID    string
	Label string
}

type victoryConditions struct {
	Standard         Victory
	LostStartingCity Victory
	GuardianArena    Victory
	HoldCity         Victory
	Tournament       Victory
}

var victoryConditionValues = victoryConditions{
	Standard:         Victory{ID: winConditions.Standard, Label: "Standard"},
	LostStartingCity: Victory{ID: winConditions.CapitalHold, Label: "Lost Starting City"},
	GuardianArena:    Victory{ID: winConditions.FinalBattle, Label: "Guardian Arena"},
	HoldCity:         Victory{ID: winConditions.CityHold, Label: "Hold City"},
	Tournament:       Victory{ID: winConditions.Tournament, Label: "Tournament"},
}

var VictoryConditions = []Victory{
	victoryConditionValues.Standard,
	victoryConditionValues.LostStartingCity,
	victoryConditionValues.GuardianArena,
	victoryConditionValues.HoldCity,
	victoryConditionValues.Tournament,
}

func GetVictoryConditionValues() victoryConditions {
	return victoryConditionValues
}

func GetVictoryCondition(id string) Victory {
	for _, victory := range VictoryConditions {
		if strings.EqualFold(victory.ID, id) {
			return victory
		}
	}
	return VictoryConditions[0] // TODO: probably should return empty Victory... suck it
}
