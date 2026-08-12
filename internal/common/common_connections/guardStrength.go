package common_connections

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

func GetPlayerToPlayerGuardStrength() models.GuardStrength {
	return models.GuardStrength{
		Default:  neutral_zone.QualityUnknown.GetGuardValue(),
		Weakest:  10000,
		Low:      22000,
		Medium:   34000,
		High:     46000,
		VeryHigh: 58000,
	}
}

func GetBronzeGuardStrength() models.GuardStrength {
	return models.GuardStrength{
		Default:  neutral_zone.QualityLow.GetGuardValue(),
		Weakest:  3000,
		Low:      6000,
		Medium:   9000,
		High:     12000,
		VeryHigh: 16000,
	}
}

func GetSilverGuardStrength() models.GuardStrength {
	return models.GuardStrength{
		Default:  neutral_zone.QualityMedium.GetGuardValue(),
		Weakest:  18000,
		Low:      21000,
		Medium:   24000,
		High:     27000,
		VeryHigh: 30000,
	}
}

func GetGoldGuardStrength() models.GuardStrength {
	return models.GuardStrength{
		Default:  neutral_zone.QualityHigh.GetGuardValue(),
		Weakest:  36000,
		Low:      42000,
		Medium:   48000,
		High:     54000,
		VeryHigh: 60000,
	}
}

func GetHubGuardStrength() models.GuardStrength {
	return models.GuardStrength{
		Default:  neutral_zone.QualityHighest.GetGuardValue(),
		Weakest:  45000,
		Low:      52000,
		Medium:   62000,
		High:     70000,
		VeryHigh: 75000,
	}
}

func GetGuardStrengthListForQuality(zoneQuality neutral_zone.Quality) []data.Tuple[string, int] {
	strength := models.GuardStrength{}
	switch zoneQuality {
	case neutral_zone.QualityLowest, neutral_zone.QualityLow:
		strength = GetBronzeGuardStrength()
	case neutral_zone.QualityMedium:
		strength = GetSilverGuardStrength()
	case neutral_zone.QualityHigh:
		strength = GetGoldGuardStrength()
	case neutral_zone.QualityHighest:
		strength = GetHubGuardStrength()
	case neutral_zone.QualityUnknown:
		fallthrough // Assume this is a player-to-player connection.
	default:
		strength = GetPlayerToPlayerGuardStrength()
	}

	return []data.Tuple[string, int]{
		data.NewTuple("Default", strength.Default),
		data.NewTuple("Weakest", strength.Weakest),
		data.NewTuple("Low", strength.Low),
		data.NewTuple("Medium", strength.Medium),
		data.NewTuple("High", strength.High),
		data.NewTuple("Very High", strength.VeryHigh),
	}
}

func GetGuardStrengthForQuality(zoneQuality neutral_zone.Quality) models.GuardStrength {
	switch zoneQuality {
	case neutral_zone.QualityLowest, neutral_zone.QualityLow:
		return GetBronzeGuardStrength()
	case neutral_zone.QualityMedium:
		return GetSilverGuardStrength()
	case neutral_zone.QualityHigh:
		return GetGoldGuardStrength()
	case neutral_zone.QualityHighest:
		return GetHubGuardStrength()
	case neutral_zone.QualityUnknown:
		fallthrough // Assume this is a player-to-player connection.
	default:
		return GetPlayerToPlayerGuardStrength()
	}
}
