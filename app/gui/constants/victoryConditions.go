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

// GetVictoryCondition looks up a victory condition by its wire ID. The second
// return value reports whether the ID is known, letting callers surface
// unknown values instead of silently reshaping them.
func GetVictoryCondition(id string) (Victory, bool) {
	for _, victory := range GetVictoryConditionList() {
		if strings.EqualFold(victory.ID, id) {
			return victory, true
		}
	}
	return Victory{}, false
}
