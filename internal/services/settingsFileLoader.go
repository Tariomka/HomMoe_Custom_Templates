package services

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// LoadSettingsFile reads a .gen.json file and returns the parsed SettingsFile.
func LoadSettingsFile(path string) (*models.EditorStateModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	settingsFile := models.NewEditorStateModel()
	if err := json.Unmarshal(data, settingsFile); err != nil {
		return nil, err
	}

	return settingsFile, nil
}

// SaveSettingsFile writes a SettingsFile to disk as indented JSON.
func SaveSettingsFile(path string, settingsFile *models.EditorStateModel) error {
	data, err := json.MarshalIndent(settingsFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// SettingsToGenerator translates a SettingsFile (UI persistence model)
// into a GeneratorSettings (generator input model).
func SettingsToGenerator(editorState *models.EditorStateModel) *generator.GeneratorSettings {
	generatorSettings := generator.NewGeneratorSettings()
	generatorSettings.TemplateName = editorState.TemplateName
	generatorSettings.GameMode = editorState.GameMode
	generatorSettings.PlayerCount = editorState.PlayerCount
	generatorSettings.MapSize = editorState.MapSize
	generatorSettings.Topology = editorState.Topology
	generatorSettings.GenerateRoads = editorState.GenerateRoads
	generatorSettings.RandomPortals = editorState.RandomPortals
	generatorSettings.SpawnRemoteFootholds = editorState.SpawnRemoteFootholds
	generatorSettings.NoDirectPlayerConnections = editorState.NoDirectPlayerConn
	generatorSettings.MaxPortalConnections = editorState.MaxPortalConnections
	generatorSettings.MinNeutralZonesBetweenPlayers = editorState.MinNeutralZonesBetweenPlayers
	generatorSettings.MatchPlayerCastleFactions = editorState.MatchPlayerCastleFactions
	generatorSettings.BannedItems = editorState.BannedItems
	generatorSettings.BannedMagics = editorState.BannedMagics
	generatorSettings.ValueOverridesText = editorState.ValueOverridesText
	generatorSettings.Bonuses = generator.ParseBonusesJSON(editorState.BonusesJSON)
	generatorSettings.PlayerZoneMandatoryContent = RowsToMandatoryContent(editorState.PlayerZoneContentRows)
	generatorSettings.LowNeutralMandatoryContent = RowsToMandatoryContent(editorState.LowNeutralContentRows)
	generatorSettings.MediumNeutralMandatoryContent = RowsToMandatoryContent(editorState.MediumNeutralContentRows)
	generatorSettings.HighNeutralMandatoryContent = RowsToMandatoryContent(editorState.HighNeutralContentRows)
	generatorSettings.HubZoneMandatoryContent = RowsToMandatoryContent(editorState.HubZoneContentRows)
	generatorSettings.FactionLawsExpPercent = editorState.FactionLawXpPercent
	generatorSettings.AstrologyExpPercent = editorState.AstrologyXpPercent

	generatorSettings.ZoneCfg = generator.ZoneConfiguration{
		NeutralZoneCount:            editorState.NeutralZoneCount,
		PlayerZoneCastles:           editorState.PlayerZoneCastles,
		NeutralZoneCastles:          editorState.NeutralZoneCastles,
		ResourceDensityPercent:      editorState.ResourceDensityPercent,
		StructureDensityPercent:     editorState.StructureDensityPercent,
		NeutralStackStrengthPercent: editorState.NeutralStackStrengthPercent,
		BorderGuardStrengthPercent:  editorState.BorderGuardStrengthPercent,
		HubZoneSize:                 editorState.HubZoneSize,
		HubZoneCastles:              editorState.HubZoneCastles,
		Advanced: generator.AdvancedSettings{
			Enabled:                    editorState.AdvancedMode,
			NeutralLowNoCastleCount:    editorState.NeutralLowNoCastleCount,
			NeutralLowCastleCount:      editorState.NeutralLowCastleCount,
			NeutralMediumNoCastleCount: editorState.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:   editorState.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:   editorState.NeutralHighNoCastleCount,
			NeutralHighCastleCount:     editorState.NeutralHighCastleCount,
			PlayerZoneSize:             editorState.PlayerZoneSize,
			NeutralZoneSize:            editorState.NeutralZoneSize,
			GuardRandomization:         editorState.GuardRandomization,
		},
	}

	generatorSettings.HeroSettings = generator.HeroSettings{
		HeroCountMin:       editorState.HeroCountMin,
		HeroCountMax:       editorState.HeroCountMax,
		HeroCountIncrement: editorState.HeroCountIncrement,
	}

	generatorSettings.GameEndConditions = &generator.GameEndConditions{
		VictoryCondition: editorState.VictoryCondition,
		CityHold:         editorState.CityHold || editorState.VictoryCondition == "win_condition_5",
		CityHoldDays:     editorState.CityHoldDays,
		LostStartCity:    editorState.LostStartCity,
		LostStartCityDay: editorState.LostStartCityDay,
		LostStartHero:    editorState.LostStartHero,
	}

	generatorSettings.GladiatorArenaRules = &generator.GladiatorArenaRules{
		Enabled:        editorState.GladiatorArena,
		DaysDelayStart: editorState.GladiatorArenaDaysDelayStart,
		CountDay:       editorState.GladiatorArenaCountDay,
	}

	generatorSettings.TournamentRules = &generator.TournamentRules{
		Enabled:            editorState.Tournament,
		FirstTournamentDay: editorState.TournamentFirstTournamentDay,
		Interval:           editorState.TournamentInterval,
		PointsToWin:        editorState.TournamentPointsToWin,
		SaveArmy:           editorState.TournamentSaveArmy,
	}

	return generatorSettings
}
