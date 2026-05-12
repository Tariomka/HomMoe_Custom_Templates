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
	gs := models.NewGeneratorSettings()
	gs.TemplateName = sf.TemplateName
	gs.PlayerCount = sf.PlayerCount
	gs.MapSize = sf.MapSize
	gs.Topology = sf.Topology
	gs.GenerateRoads = sf.GenerateRoads
	gs.RandomPortals = sf.RandomPortals
	gs.SpawnRemoteFootholds = sf.SpawnRemoteFootholds
	gs.NoDirectPlayerConnections = sf.NoDirectPlayerConn
	gs.MaxPortalConnections = sf.MaxPortalConnections
	gs.MinNeutralZonesBetweenPlayers = sf.MinNeutralZonesBetweenPlayers
	gs.MatchPlayerCastleFactions = sf.MatchPlayerCastleFactions
	gs.ExperimentalBalancedZonePlacement = sf.ExperimentalBalancedZonePlacement
	gs.FactionLawsExpPercent = sf.FactionLawsExpPercent
	gs.AstrologyExpPercent = sf.AstrologyExpPercent

	gs.ZoneCfg = models.ZoneConfiguration{
		NeutralZoneCount:            sf.NeutralZoneCount,
		PlayerZoneCastles:           sf.PlayerZoneCastles,
		NeutralZoneCastles:          sf.NeutralZoneCastles,
		ResourceDensityPercent:      sf.EffectiveResourceDensity(),
		StructureDensityPercent:     sf.EffectiveStructureDensity(),
		NeutralStackStrengthPercent: sf.NeutralStackStrengthPercent,
		BorderGuardStrengthPercent:  sf.BorderGuardStrengthPercent,
		HubZoneSize:                 sf.HubZoneSize,
		HubZoneCastles:              sf.HubZoneCastles,
		Advanced: models.AdvancedSettings{
			Enabled:                    sf.AdvancedMode,
			NeutralLowNoCastleCount:    sf.NeutralLowNoCastleCount,
			NeutralLowCastleCount:      sf.NeutralLowCastleCount,
			NeutralMediumNoCastleCount: sf.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:   sf.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:   sf.NeutralHighNoCastleCount,
			NeutralHighCastleCount:     sf.NeutralHighCastleCount,
			PlayerZoneSize:             sf.PlayerZoneSize,
			NeutralZoneSize:            sf.NeutralZoneSize,
			GuardRandomization:         sf.GuardRandomization,
		},
	}

	gs.HeroSettings = &models.HeroSettings{
		HeroCountMin:       sf.HeroCountMin,
		HeroCountMax:       sf.HeroCountMax,
		HeroCountIncrement: sf.HeroCountIncrement,
	}

	gs.GameEndConditions = &models.GameEndConditions{
		VictoryCondition: sf.VictoryCondition,
		CityHold:         sf.CityHold || sf.VictoryCondition == "win_condition_5",
		CityHoldDays:     sf.CityHoldDays,
		LostStartCity:    sf.LostStartCity,
		LostStartCityDay: sf.LostStartCityDay,
		LostStartHero:    sf.LostStartHero,
	}

	gs.GladiatorArenaRules = &models.GladiatorArenaRules{
		Enabled:        sf.GladiatorArena,
		DaysDelayStart: sf.GladiatorArenaDaysDelayStart,
		CountDay:       sf.GladiatorArenaCountDay,
	}

	gs.TournamentRules = &models.TournamentRules{
		Enabled:            sf.Tournament,
		FirstTournamentDay: sf.TournamentFirstTournamentDay,
		Interval:           sf.TournamentInterval,
		PointsToWin:        sf.TournamentPointsToWin,
		SaveArmy:           sf.TournamentSaveArmy,
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
