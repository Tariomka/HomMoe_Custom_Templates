package services

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// LoadSettingsFile reads a .gen.json file and returns the parsed SettingsFile.
func LoadSettingsFile(path string) (*models.SettingsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	settingsFile := models.NewSettingsFile()
	if err := json.Unmarshal(data, settingsFile); err != nil {
		return nil, err
	}
	// One-way legacy upgrade: older .gen.json files used a boolean flag
	// for the balanced topology. Promote it to Topology = Balanced and
	// clear the flag so subsequent saves drop it from disk.
	if settingsFile.ExperimentalBalancedZonePlacement {
		settingsFile.Topology = generator.TopologyBalanced
		settingsFile.ExperimentalBalancedZonePlacement = false
	}
	return settingsFile, nil
}

// SaveSettingsFile writes a SettingsFile to disk as indented JSON.
func SaveSettingsFile(path string, settingsFile *models.SettingsFile) error {
	data, err := json.MarshalIndent(settingsFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SettingsToGenerator translates a SettingsFile (UI persistence model)
// into a GeneratorSettings (generator input model).
func SettingsToGenerator(settingsFile *models.SettingsFile) *models.GeneratorSettings {
	generatorSettings := generator.NewGeneratorSettings()
	generatorSettings.TemplateName = settingsFile.TemplateName
	generatorSettings.PlayerCount = settingsFile.PlayerCount
	generatorSettings.MapSize = settingsFile.MapSize
	generatorSettings.Topology = settingsFile.Topology
	generatorSettings.GenerateRoads = settingsFile.GenerateRoads
	generatorSettings.RandomPortals = settingsFile.RandomPortals
	generatorSettings.SpawnRemoteFootholds = settingsFile.SpawnRemoteFootholds
	generatorSettings.NoDirectPlayerConnections = settingsFile.NoDirectPlayerConn
	generatorSettings.MaxPortalConnections = settingsFile.MaxPortalConnections
	generatorSettings.MinNeutralZonesBetweenPlayers = settingsFile.MinNeutralZonesBetweenPlayers
	generatorSettings.MatchPlayerCastleFactions = settingsFile.MatchPlayerCastleFactions
	generatorSettings.BannedItems = settingsFile.BannedItems
	generatorSettings.BannedMagics = settingsFile.BannedMagics
	generatorSettings.ValueOverridesText = settingsFile.ValueOverridesText
	generatorSettings.Bonuses = generator.ParseBonusesJSON(settingsFile.BonusesJson)
	generatorSettings.PlayerZoneMandatoryContent = RowsToMandatoryContent(settingsFile.PlayerZoneContentRows)
	generatorSettings.LowNeutralMandatoryContent = RowsToMandatoryContent(settingsFile.LowNeutralContentRows)
	generatorSettings.MediumNeutralMandatoryContent = RowsToMandatoryContent(settingsFile.MediumNeutralContentRows)
	generatorSettings.HighNeutralMandatoryContent = RowsToMandatoryContent(settingsFile.HighNeutralContentRows)
	generatorSettings.HubZoneMandatoryContent = RowsToMandatoryContent(settingsFile.HubZoneContentRows)
	generatorSettings.FactionLawsExpPercent = settingsFile.FactionLawsExpPercent
	generatorSettings.AstrologyExpPercent = settingsFile.AstrologyExpPercent

	generatorSettings.ZoneCfg = models.ZoneConfiguration{
		NeutralZoneCount:            settingsFile.NeutralZoneCount,
		PlayerZoneCastles:           settingsFile.PlayerZoneCastles,
		NeutralZoneCastles:          settingsFile.NeutralZoneCastles,
		ResourceDensityPercent:      settingsFile.EffectiveResourceDensity(),
		StructureDensityPercent:     settingsFile.EffectiveStructureDensity(),
		NeutralStackStrengthPercent: settingsFile.NeutralStackStrengthPercent,
		BorderGuardStrengthPercent:  settingsFile.BorderGuardStrengthPercent,
		HubZoneSize:                 settingsFile.HubZoneSize,
		HubZoneCastles:              settingsFile.HubZoneCastles,
		Advanced: models.AdvancedSettings{
			Enabled:                    settingsFile.AdvancedMode,
			NeutralLowNoCastleCount:    settingsFile.NeutralLowNoCastleCount,
			NeutralLowCastleCount:      settingsFile.NeutralLowCastleCount,
			NeutralMediumNoCastleCount: settingsFile.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:   settingsFile.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:   settingsFile.NeutralHighNoCastleCount,
			NeutralHighCastleCount:     settingsFile.NeutralHighCastleCount,
			PlayerZoneSize:             settingsFile.PlayerZoneSize,
			NeutralZoneSize:            settingsFile.NeutralZoneSize,
			GuardRandomization:         settingsFile.GuardRandomization,
		},
	}

	generatorSettings.HeroSettings = &models.HeroSettings{
		HeroCountMin:       settingsFile.HeroCountMin,
		HeroCountMax:       settingsFile.HeroCountMax,
		HeroCountIncrement: settingsFile.HeroCountIncrement,
	}

	generatorSettings.GameEndConditions = &models.GameEndConditions{
		VictoryCondition: settingsFile.VictoryCondition,
		CityHold:         settingsFile.CityHold || settingsFile.VictoryCondition == "win_condition_5",
		CityHoldDays:     settingsFile.CityHoldDays,
		LostStartCity:    settingsFile.LostStartCity,
		LostStartCityDay: settingsFile.LostStartCityDay,
		LostStartHero:    settingsFile.LostStartHero,
	}

	generatorSettings.GladiatorArenaRules = &models.GladiatorArenaRules{
		Enabled:        settingsFile.GladiatorArena,
		DaysDelayStart: settingsFile.GladiatorArenaDaysDelayStart,
		CountDay:       settingsFile.GladiatorArenaCountDay,
	}

	generatorSettings.TournamentRules = &models.TournamentRules{
		Enabled:            settingsFile.Tournament,
		FirstTournamentDay: settingsFile.TournamentFirstTournamentDay,
		Interval:           settingsFile.TournamentInterval,
		PointsToWin:        settingsFile.TournamentPointsToWin,
		SaveArmy:           settingsFile.TournamentSaveArmy,
	}

	return generatorSettings
}
