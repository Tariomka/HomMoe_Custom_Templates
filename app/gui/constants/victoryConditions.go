package constants

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

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

func GetVictoryConditionList() []Victory {
	values := GetVictoryConditionValues()
	return []Victory{
		values.Standard,
		values.LostStartingCity,
		values.GuardianArena,
		values.HoldCity,
		values.Tournament,
	}
}

func GetVictoryConditionValues() victoryConditions {
	winConditions := registry.GetWinningConditionValues()
	return victoryConditions{
		Standard:         Victory{ID: winConditions.Standard, Label: "Standard"},
		LostStartingCity: Victory{ID: winConditions.CapitalHold, Label: "Lost Starting City"},
		GuardianArena:    Victory{ID: winConditions.FinalBattle, Label: "Guardian Arena"},
		HoldCity:         Victory{ID: winConditions.CityHold, Label: "Hold City"},
		Tournament:       Victory{ID: winConditions.Tournament, Label: "Tournament"},
	}
}

func GetVictoryCondition(id string) Victory {
	for _, victory := range GetVictoryConditionList() {
		if strings.EqualFold(victory.ID, id) {
			return victory
		}
	}
	return GetVictoryConditionList()[0] // TODO: probably should return empty Victory... suck it
}
