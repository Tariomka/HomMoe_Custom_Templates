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
	generatorSettings.Bonuses = config.ParseBonusesJSON(editorState.BonusesJSON)
	generatorSettings.PlayerZoneMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.PlayerZoneContentRows)
	generatorSettings.LowNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.LowNeutralContentRows)
	generatorSettings.MediumNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.MediumNeutralContentRows)
	generatorSettings.HighNeutralMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.HighNeutralContentRows)
	generatorSettings.HubZoneMandatoryContent = contentProvider.CreateContentItemsFrom(editorState.HubZoneContentRows)
	generatorSettings.FactionLawsExpPercent = editorState.FactionLawXpPercent
	generatorSettings.AstrologyExpPercent = editorState.AstrologyXpPercent

	generatorSettings.ZoneConfiguration = config.ZoneConfig{
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

	generatorSettings.HeroSettings = config.HeroSettings{
		HeroCountMin:       editorState.HeroCountMin,
		HeroCountMax:       editorState.HeroCountMax,
		HeroCountIncrement: editorState.HeroCountIncrement,
	}

	generatorSettings.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: editorState.VictoryCondition,
		CityHold:         editorState.CityHold || editorState.VictoryCondition == winConditions.CityHold,
		CityHoldDays:     editorState.CityHoldDays,
		LostStartCity:    editorState.LostStartCity,
		LostStartCityDay: editorState.LostStartCityDay,
		LostStartHero:    editorState.LostStartHero,
	}

	generatorSettings.GladiatorArenaRules = &config.GladiatorArenaRules{
		Enabled:        editorState.GladiatorArena,
		DaysDelayStart: editorState.GladiatorArenaDaysDelayStart,
		CountDay:       editorState.GladiatorArenaCountDay,
	}

	generatorSettings.TournamentRules = &config.TournamentRules{
		Enabled:            editorState.Tournament,
		FirstTournamentDay: editorState.TournamentFirstTournamentDay,
		Interval:           editorState.TournamentInterval,
		PointsToWin:        editorState.TournamentPointsToWin,
		SaveArmy:           editorState.TournamentSaveArmy,
	}

	return generatorSettings
}
