package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
)

var (
	winConditions = registry.GetWinningConditionValues()
)

type GeneratorConfigMapper struct{}

func NewConfigMapper() *GeneratorConfigMapper {
	return &GeneratorConfigMapper{}
}

// FromEditorState translates a SettingsFile (UI persistence model)
// into a GeneratorSettings (generator input model).
func (this *GeneratorConfigMapper) FromEditorState(editorState dtos.EditorStateDto) *config.GeneratorConfig {
	contentProvider := providers.NewMandatoryContentProvider()
	generatorSettings := config.NewGeneratorConfig()
	generatorSettings.TemplateName = editorState.TemplateName
	generatorSettings.GameMode = editorState.GameMode
	generatorSettings.PlayerCount = editorState.PlayerCount
	generatorSettings.MapSize = editorState.MapSize
	generatorSettings.Topology = editorState.Topology
	generatorSettings.GenerateRoads = editorState.GenerateRoads
	generatorSettings.RandomPortals = editorState.RandomPortals
	generatorSettings.SpawnRemoteFootholds = editorState.SpawnRemoteFootholds
	generatorSettings.RemoteFootholdCount = editorState.RemoteFootholdCount
	generatorSettings.NoDirectPlayerConnections = editorState.NoDirectPlayerConn
	generatorSettings.MaxPortalConnections = editorState.MaxPortalConnections
	generatorSettings.MinNeutralZonesBetweenPlayers = editorState.MinNeutralZonesBetweenPlayers
	generatorSettings.MatchPlayerCastleFactions = editorState.MatchPlayerCastleFactions
	generatorSettings.BannedItems = editorState.BannedItems
	generatorSettings.BannedMagics = editorState.BannedMagics
	generatorSettings.ValueOverridesText = editorState.ValueOverridesText
	generatorSettings.Bonuses = editorState.Bonuses
	generatorSettings.PlayerZoneMandatoryContent = contentProvider.CreateContentItemsFrom(
		editorState.PlayerZoneContentRows)
	generatorSettings.LowNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(
		editorState.LowNeutralContentRows)
	generatorSettings.MediumNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(
		editorState.MediumNeutralContentRows)
	generatorSettings.HighNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(
		editorState.HighNeutralContentRows)
	generatorSettings.HubZoneMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.HubZoneContentRows)
	generatorSettings.FactionLawsExpPercent = editorState.FactionLawXpPercent
	generatorSettings.AstrologyExpPercent = editorState.AstrologyXpPercent

	generatorSettings.ZoneConfiguration = this.mapZoneConfig(editorState)
	generatorSettings.HeroSettings = this.mapHeroSettings(editorState)
	generatorSettings.GameEndConditions = this.mapGameEndConditions(editorState)
	generatorSettings.GladiatorArenaRules = this.mapGladiatorArenaRules(editorState)
	generatorSettings.TournamentRules = this.mapTournamentRules(editorState)

	return generatorSettings
}

func (this *GeneratorConfigMapper) mapZoneConfig(editorState dtos.EditorStateDto) config.ZoneConfig {
	return config.ZoneConfig{
		NeutralZoneCount:            editorState.NeutralZoneCount,
		PlayerOwnedCastles:          editorState.PlayerOwnedCastles,
		PlayerZoneCastles:           editorState.PlayerZoneCastles,
		NeutralZoneCastles:          editorState.NeutralZoneCastles,
		SpawnAbandonedOutposts:      editorState.SpawnAbandonedOutposts,
		AbandonedOutpostCount:       editorState.AbandonedOutpostCount,
		ResourceDensityPercent:      editorState.ResourceDensityPercent,
		StructureDensityPercent:     editorState.StructureDensityPercent,
		NeutralStackStrengthPercent: editorState.NeutralStackStrengthPercent,
		BorderGuardStrengthPercent:  editorState.BorderGuardStrengthPercent,
		HubZoneSize:                 editorState.HubZoneSize,
		HubZoneCastles:              editorState.HubZoneCastles,
		Advanced: config.AdvancedSettings{
			Enabled:                     editorState.AdvancedMode,
			NeutralLowNoCastleCount:     editorState.NeutralLowNoCastleCount,
			NeutralLowCastleCount:       editorState.NeutralLowCastleCount,
			NeutralMediumNoCastleCount:  editorState.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:    editorState.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:    editorState.NeutralHighNoCastleCount,
			NeutralHighCastleCount:      editorState.NeutralHighCastleCount,
			NeutralLowCastlesPerZone:    editorState.NeutralLowCastlesPerZone,
			NeutralMediumCastlesPerZone: editorState.NeutralMediumCastlesPerZone,
			NeutralHighCastlesPerZone:   editorState.NeutralHighCastlesPerZone,
			PlayerZoneSize:              editorState.PlayerZoneSize,
			NeutralZoneSize:             editorState.NeutralZoneSize,
			GuardRandomization:          editorState.GuardRandomization,
		},
	}
}

func (this *GeneratorConfigMapper) mapHeroSettings(editorState dtos.EditorStateDto) config.HeroSettings {
	return config.HeroSettings{
		HeroCountMin:       editorState.HeroCountMin,
		HeroCountMax:       editorState.HeroCountMax,
		HeroCountIncrement: editorState.HeroCountIncrement,
	}
}

func (this *GeneratorConfigMapper) mapGameEndConditions(editorState dtos.EditorStateDto) *config.GameEndConditions {
	return &config.GameEndConditions{
		VictoryCondition: editorState.VictoryCondition,
		CityHold:         editorState.CityHold || editorState.VictoryCondition == winConditions.CityHold,
		CityHoldDays:     editorState.CityHoldDays,
		LostStartCity:    editorState.LostStartCity,
		LostStartCityDay: editorState.LostStartCityDay,
		LostStartHero:    editorState.LostStartHero,
	}
}

func (this *GeneratorConfigMapper) mapGladiatorArenaRules(editorState dtos.EditorStateDto) *config.GladiatorArenaRules {
	return &config.GladiatorArenaRules{
		Enabled:        editorState.GladiatorArena,
		DaysDelayStart: editorState.GladiatorArenaDaysDelayStart,
		CountDay:       editorState.GladiatorArenaCountDay,
	}
}

func (this *GeneratorConfigMapper) mapTournamentRules(editorState dtos.EditorStateDto) *config.TournamentRules {
	return &config.TournamentRules{
		Enabled:            editorState.Tournament,
		FirstTournamentDay: editorState.TournamentFirstTournamentDay,
		Interval:           editorState.TournamentInterval,
		PointsToWin:        editorState.TournamentPointsToWin,
		SaveArmy:           editorState.TournamentSaveArmy,
	}
}
