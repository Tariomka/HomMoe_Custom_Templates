package gui

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// LoadSettingsFile reads a .oetgs file and returns the parsed SettingsFile.
func LoadSettingsFile(path string) (*models.SettingsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sf := models.NewSettingsFile()
	if err := json.Unmarshal(data, sf); err != nil {
		return nil, err
	}
	return sf, nil
}

// SaveSettingsFile writes a SettingsFile to disk as indented JSON.
func SaveSettingsFile(path string, sf *models.SettingsFile) error {
	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// settingsFileToGenerator translates a SettingsFile (UI persistence model)
// into a GeneratorSettings (generator input model).
func settingsFileToGenerator(sf *models.SettingsFile) *models.GeneratorSettings {
	gs := &models.GeneratorSettings{
		TemplateName:                  sf.TemplateName,
		GameMode:                      "Classic",
		PlayerCount:                   sf.PlayerCount,
		MapSize:                       mapSizeLabelInt(sf.MapSize),
		MapSizeValue:                  sf.MapSize,
		Topology:                      sf.Topology,
		AllowRoads:                    sf.GenerateRoads,
		AllowPortals:                  sf.RandomPortals,
		AllowFootholds:                sf.SpawnRemoteFootholds,
		EnablePlayerIsolation:         sf.NoDirectPlayerConn,
		EnableCityHold:                sf.CityHold,
		ShowDescription:               true,
		IncludeOptionsInDescription:   true,
		MaxPortalConnections:          sf.MaxPortalConnections,
		MinNeutralZonesBetweenPlayers: sf.MinNeutralZonesBetweenPlayers,
		MatchPlayerCastleFactions:     sf.MatchPlayerCastleFactions,
		NeutralStackStrengthPercent:   sf.NeutralStackStrengthPercent,
		BorderGuardStrengthPercent:    sf.BorderGuardStrengthPercent,
		ResourceDensityPercent:        sf.EffectiveResourceDensity(),
		StructureDensityPercent:       sf.EffectiveStructureDensity(),
		FactionLawsExpPercent:         sf.FactionLawsExpPercent,
		AstrologyExpPercent:           sf.AstrologyExpPercent,
		PlayerZoneSize:                sf.PlayerZoneSize,
		NeutralZoneSize:               sf.NeutralZoneSize,
		HubZoneSize:                   sf.HubZoneSize,
		HubZoneCastles:                sf.HubZoneCastles,
		PlayerZoneMandatoryContent:    sf.PlayerZoneMandatoryContent,
		AdvancedSettings: &models.AdvancedSettings{
			Enabled:                    sf.AdvancedMode,
			GuardRandomization:         sf.GuardRandomization,
			ConnectionCountPerZone:     2,
			NeutralZoneCount:           sf.NeutralZoneCount,
			NeutralLowNoCastleCount:    sf.NeutralLowNoCastleCount,
			NeutralLowCastleCount:      sf.NeutralLowCastleCount,
			NeutralMediumNoCastleCount: sf.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:   sf.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:   sf.NeutralHighNoCastleCount,
			NeutralHighCastleCount:     sf.NeutralHighCastleCount,
			PlayerCastleCount:          sf.PlayerZoneCastles,
			NeutralCastleCount:         sf.NeutralZoneCastles,
			NeutralZoneLowCount:        sf.NeutralLowNoCastleCount + sf.NeutralLowCastleCount,
			NeutralZoneMediumCount:     sf.NeutralMediumNoCastleCount + sf.NeutralMediumCastleCount,
			NeutralZoneHighCount:       sf.NeutralHighNoCastleCount + sf.NeutralHighCastleCount,
		},
		HeroSettings: &models.HeroSettings{
			HeroCount:          sf.HeroCountMax,
			HeroCountMin:       sf.HeroCountMin,
			HeroCountMax:       sf.HeroCountMax,
			HeroCountIncrement: sf.HeroCountIncrement,
		},
		GameEndConditions: &models.GameEndConditions{
			VictoryCondition:     sf.VictoryCondition,
			EnableClassicVictory: sf.VictoryCondition == "win_condition_1",
			EnableCityHold:       sf.CityHold || sf.VictoryCondition == "win_condition_5",
			CityHoldDays:         sf.CityHoldDays,
			LostStartCity:        sf.LostStartCity,
			LostStartCityDay:     sf.LostStartCityDay,
			LostStartHero:        sf.LostStartHero,
			EnableGladiatorArena: sf.GladiatorArena,
			EnableTournaments:    sf.Tournament || sf.VictoryCondition == "win_condition_6",
		},
		GladiatorArenaRules: &models.GladiatorArenaRules{
			Enabled:        sf.GladiatorArena,
			DaysDelayStart: sf.GladiatorArenaDaysDelayStart,
			CountDay:       sf.GladiatorArenaCountDay,
		},
		TournamentRules: &models.TournamentRules{
			Enabled:            sf.Tournament,
			FirstTournamentDay: sf.TournamentFirstTournamentDay,
			Interval:           sf.TournamentInterval,
			PointsToWin:        sf.TournamentPointsToWin,
			SaveArmy:           sf.TournamentSaveArmy,
		},
	}
	return gs
}

// mapSizeLabelInt returns the short S/M/L/XL/H/G/C label for an integer size.
// Mirrors KnownValues.MapSizeLabel from the C# source.
func mapSizeLabelInt(size int) string {
	switch {
	case size == 64:
		return "S"
	case size == 80 || size == 96:
		return "M"
	case size == 112 || size == 128:
		return "L"
	case size == 144 || size == 160:
		return "XL"
	case size == 176 || size == 192:
		return "H"
	case size >= 208 && size <= 256:
		return "G"
	default:
		return "C"
	}
}
