package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type GeneratorConfigMapper struct {
	contentItemMapper *MandatoryContentItemMapper
}

func NewConfigMapper(contentItemMapper *MandatoryContentItemMapper) *GeneratorConfigMapper {
	return &GeneratorConfigMapper{
		contentItemMapper: contentItemMapper,
	}
}

// FromEditorState translates a SettingsFile (UI persistence model)
// into a GeneratorSettings (generator input model).
func (this *GeneratorConfigMapper) FromEditorState(editorState dtos.EditorStateDto) *config.GeneratorConfig {
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
	generatorSettings.MatchPlayerCastleFactions = editorState.MatchPlayerCastleFactions
	generatorSettings.BannedItems = editorState.BannedItems
	generatorSettings.BannedMagics = editorState.BannedMagics
	generatorSettings.ValueOverridesText = editorState.ValueOverridesText
	generatorSettings.Bonuses = editorState.Bonuses
	generatorSettings.PlayerZoneMandatoryContent = this.contentItemMapper.FromRows(
		editorState.PlayerZoneContentRows)
	generatorSettings.LowestNeutralMandatoryContent = this.contentItemMapper.FromRows(
		editorState.LowestNeutralContentRows)
	generatorSettings.LowNeutralMandatoryContent = this.contentItemMapper.FromRows(
		editorState.LowNeutralContentRows)
	generatorSettings.MediumNeutralMandatoryContent = this.contentItemMapper.FromRows(
		editorState.MediumNeutralContentRows)
	generatorSettings.HighNeutralMandatoryContent = this.contentItemMapper.FromRows(
		editorState.HighNeutralContentRows)
	generatorSettings.HubZoneMandatoryContent = this.contentItemMapper.FromRows(
		editorState.HubZoneContentRows)
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
		PlayerZoneSize:              editorState.PlayerZoneSize,
		NeutralZoneSize:             editorState.NeutralZoneSize,
		GuardRandomization:          editorState.GuardRandomization,
		Advanced: config.AdvancedSettings{
			Enabled:                     editorState.AdvancedMode,
			NeutralLowestNoCastleCount:  editorState.NeutralLowestNoCastleCount,
			NeutralLowestCastleCount:    editorState.NeutralLowestCastleCount,
			NeutralLowNoCastleCount:     editorState.NeutralLowNoCastleCount,
			NeutralLowCastleCount:       editorState.NeutralLowCastleCount,
			NeutralMediumNoCastleCount:  editorState.NeutralMediumNoCastleCount,
			NeutralMediumCastleCount:    editorState.NeutralMediumCastleCount,
			NeutralHighNoCastleCount:    editorState.NeutralHighNoCastleCount,
			NeutralHighCastleCount:      editorState.NeutralHighCastleCount,
			NeutralLowestCastlesPerZone: editorState.NeutralLowestCastlesPerZone,
			NeutralLowCastlesPerZone:    editorState.NeutralLowCastlesPerZone,
			NeutralMediumCastlesPerZone: editorState.NeutralMediumCastlesPerZone,
			NeutralHighCastlesPerZone:   editorState.NeutralHighCastlesPerZone,
			HubZoneCastles:              editorState.HubZoneCastles,
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
		CityHold:         editorState.CityHold,
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
