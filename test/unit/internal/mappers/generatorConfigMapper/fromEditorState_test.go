package generatorConfigMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenScalarOptionsProvided_CopiesEachToConfig(t *testing.T) {
	templateName := gofakeit.ProductName()
	bannedItems := gofakeit.Word()
	bannedMagics := gofakeit.Word()
	valueOverrides := gofakeit.Word()
	playerCount := gofakeit.Number(2, 8)
	mapSize := gofakeit.Number(96, 224)
	maxPortals := gofakeit.Number(1, 16)
	minNeutralsBetween := gofakeit.Number(1, 4)
	remoteFootholds := gofakeit.Number(1, 4)
	factionLawXp := gofakeit.Number(50, 200)
	astrologyXp := gofakeit.Number(50, 200)

	testCases := []struct {
		subtestName string
		mutate      func(state *dtos.EditorStateDto)
		actual      func(configuration *config.GeneratorConfig) any
		expected    any
	}{
		{
			"WhenTemplateNameProvided_CopiesTemplateName",
			func(state *dtos.EditorStateDto) { state.TemplateName = templateName },
			func(configuration *config.GeneratorConfig) any { return configuration.TemplateName },
			templateName,
		},
		{
			"WhenGameModeProvided_CopiesGameMode",
			func(state *dtos.EditorStateDto) { state.GameMode = "SingleHero" },
			func(configuration *config.GeneratorConfig) any { return configuration.GameMode },
			"SingleHero",
		},
		{
			"WhenPlayerCountProvided_CopiesPlayerCount",
			func(state *dtos.EditorStateDto) { state.PlayerCount = playerCount },
			func(configuration *config.GeneratorConfig) any { return configuration.PlayerCount },
			playerCount,
		},
		{
			"WhenMapSizeProvided_CopiesMapSize",
			func(state *dtos.EditorStateDto) { state.MapSize = mapSize },
			func(configuration *config.GeneratorConfig) any { return configuration.MapSize },
			mapSize,
		},
		{
			"WhenTopologyProvided_CopiesTopology",
			func(state *dtos.EditorStateDto) { state.Topology = config.TopologyChain },
			func(configuration *config.GeneratorConfig) any { return configuration.Topology },
			config.TopologyChain,
		},
		{
			"WhenGenerateRoadsEnabled_CopiesFlag",
			func(state *dtos.EditorStateDto) { state.GenerateRoads = true },
			func(configuration *config.GeneratorConfig) any { return configuration.GenerateRoads },
			true,
		},
		{
			"WhenRandomPortalsEnabled_CopiesFlag",
			func(state *dtos.EditorStateDto) { state.RandomPortals = true },
			func(configuration *config.GeneratorConfig) any { return configuration.RandomPortals },
			true,
		},
		{
			"WhenMaxPortalConnectionsProvided_CopiesCount",
			func(state *dtos.EditorStateDto) { state.MaxPortalConnections = maxPortals },
			func(configuration *config.GeneratorConfig) any { return configuration.MaxPortalConnections },
			maxPortals,
		},
		{
			"WhenMinNeutralZonesBetweenPlayersProvided_CopiesCount",
			func(state *dtos.EditorStateDto) { state.MinNeutralZonesBetweenPlayers = minNeutralsBetween },
			func(configuration *config.GeneratorConfig) any { return configuration.MinNeutralZonesBetweenPlayers },
			minNeutralsBetween,
		},
		{
			"WhenMatchPlayerCastleFactionsEnabled_CopiesFlag",
			func(state *dtos.EditorStateDto) { state.MatchPlayerCastleFactions = true },
			func(configuration *config.GeneratorConfig) any { return configuration.MatchPlayerCastleFactions },
			true,
		},
		{
			"WhenSpawnRemoteFootholdsDisabled_CopiesFlag",
			func(state *dtos.EditorStateDto) { state.SpawnRemoteFootholds = false },
			func(configuration *config.GeneratorConfig) any { return configuration.SpawnRemoteFootholds },
			false,
		},
		{
			"WhenRemoteFootholdCountProvided_CopiesCount",
			func(state *dtos.EditorStateDto) { state.RemoteFootholdCount = remoteFootholds },
			func(configuration *config.GeneratorConfig) any { return configuration.RemoteFootholdCount },
			remoteFootholds,
		},
		{
			"WhenNoDirectPlayerConnEnabled_CopiesFlagToNoDirectPlayerConnections",
			func(state *dtos.EditorStateDto) { state.NoDirectPlayerConn = true },
			func(configuration *config.GeneratorConfig) any { return configuration.NoDirectPlayerConnections },
			true,
		},
		{
			"WhenBannedItemsProvided_CopiesText",
			func(state *dtos.EditorStateDto) { state.BannedItems = bannedItems },
			func(configuration *config.GeneratorConfig) any { return configuration.BannedItems },
			bannedItems,
		},
		{
			"WhenBannedMagicsProvided_CopiesText",
			func(state *dtos.EditorStateDto) { state.BannedMagics = bannedMagics },
			func(configuration *config.GeneratorConfig) any { return configuration.BannedMagics },
			bannedMagics,
		},
		{
			"WhenValueOverridesTextProvided_CopiesText",
			func(state *dtos.EditorStateDto) { state.ValueOverridesText = valueOverrides },
			func(configuration *config.GeneratorConfig) any { return configuration.ValueOverridesText },
			valueOverrides,
		},
		{
			"WhenFactionLawXpPercentProvided_CopiesPercent",
			func(state *dtos.EditorStateDto) { state.FactionLawXpPercent = factionLawXp },
			func(configuration *config.GeneratorConfig) any { return configuration.FactionLawsExpPercent },
			factionLawXp,
		},
		{
			"WhenAstrologyXpPercentProvided_CopiesPercent",
			func(state *dtos.EditorStateDto) { state.AstrologyXpPercent = astrologyXp },
			func(configuration *config.GeneratorConfig) any { return configuration.AstrologyExpPercent },
			astrologyXp,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			state := dtos.NewDefaultEditorStateDto()
			testCase.mutate(&state)

			// Act
			configuration := mappers.NewConfigMapper().FromEditorState(state)

			// Assert
			assert.Equal(t, testCase.expected, testCase.actual(configuration))
		})
	}
}

func TestWhenZoneOptionsProvided_PopulatesZoneConfiguration(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.NeutralZoneCount = 5
	state.PlayerOwnedCastles = 1
	state.PlayerZoneCastles = 2
	state.NeutralZoneCastles = 3
	state.SpawnAbandonedOutposts = true
	state.AbandonedOutpostCount = 2
	state.ResourceDensityPercent = 150
	state.StructureDensityPercent = 80
	state.NeutralStackStrengthPercent = 110
	state.BorderGuardStrengthPercent = 90
	state.HubZoneSize = 1.5
	state.HubZoneCastles = 1
	state.AdvancedMode = true
	state.NeutralLowNoCastleCount = 1
	state.NeutralLowCastleCount = 2
	state.NeutralMediumNoCastleCount = 3
	state.NeutralMediumCastleCount = 2
	state.NeutralHighNoCastleCount = 4
	state.NeutralHighCastleCount = 1
	state.NeutralLowCastlesPerZone = 1
	state.NeutralMediumCastlesPerZone = 2
	state.NeutralHighCastlesPerZone = 3
	state.PlayerZoneSize = 1.2
	state.NeutralZoneSize = 0.8
	state.GuardRandomization = 0.1
	expected := config.ZoneConfig{
		NeutralZoneCount:            5,
		PlayerOwnedCastles:          1,
		PlayerZoneCastles:           2,
		NeutralZoneCastles:          3,
		SpawnAbandonedOutposts:      true,
		AbandonedOutpostCount:       2,
		ResourceDensityPercent:      150,
		StructureDensityPercent:     80,
		NeutralStackStrengthPercent: 110,
		BorderGuardStrengthPercent:  90,
		HubZoneSize:                 1.5,
		HubZoneCastles:              1,
		Advanced: config.AdvancedSettings{
			Enabled:                     true,
			NeutralLowNoCastleCount:     1,
			NeutralLowCastleCount:       2,
			NeutralMediumNoCastleCount:  3,
			NeutralMediumCastleCount:    2,
			NeutralHighNoCastleCount:    4,
			NeutralHighCastleCount:      1,
			NeutralLowCastlesPerZone:    1,
			NeutralMediumCastlesPerZone: 2,
			NeutralHighCastlesPerZone:   3,
			PlayerZoneSize:              1.2,
			NeutralZoneSize:             0.8,
			GuardRandomization:          0.1,
		},
	}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, expected, configuration.ZoneConfiguration)
}

func TestWhenHeroOptionsProvided_PopulatesHeroSettings(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.HeroCountMin = 2
	state.HeroCountMax = 9
	state.HeroCountIncrement = 3
	expected := config.HeroSettings{HeroCountMin: 2, HeroCountMax: 9, HeroCountIncrement: 3}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, expected, configuration.HeroSettings)
}

func TestWhenManualCityHoldOptionsProvided_PopulatesGameEndConditions(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.VictoryCondition = "win_condition_1"
	state.CityHold = true
	state.CityHoldDays = 10
	state.LostStartCity = true
	state.LostStartCityDay = 4
	state.LostStartHero = true
	expected := &config.GameEndConditions{
		VictoryCondition: "win_condition_1",
		CityHold:         true,
		CityHoldDays:     10,
		LostStartCity:    true,
		LostStartCityDay: 4,
		LostStartHero:    true,
	}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, expected, configuration.GameEndConditions)
}

func TestWhenVictoryConditionIsCityHoldCondition_ForcesCityHoldEnabled(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.VictoryCondition = "win_condition_5"
	state.CityHold = false

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.True(t, configuration.GameEndConditions.CityHold)
}

func TestWhenGladiatorArenaOptionsProvided_PopulatesGladiatorArenaRules(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.GladiatorArena = true
	state.GladiatorArenaDaysDelayStart = 12
	state.GladiatorArenaCountDay = 4
	expected := &config.GladiatorArenaRules{Enabled: true, DaysDelayStart: 12, CountDay: 4}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, expected, configuration.GladiatorArenaRules)
}

func TestWhenTournamentOptionsProvided_PopulatesTournamentRules(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.Tournament = true
	state.TournamentFirstTournamentDay = 21
	state.TournamentInterval = 5
	state.TournamentPointsToWin = 4
	state.TournamentSaveArmy = true
	expected := &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 21,
		Interval:           5,
		PointsToWin:        4,
		SaveArmy:           true,
	}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, expected, configuration.TournamentRules)
}

func TestWhenContentRowsProvidedForEveryZoneKind_PopulatesEveryMandatoryCollection(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "a", Count: 1}}
	state.LowNeutralContentRows = []models.ZoneContentRowSave{{Sid: "b", Count: 1}}
	state.MediumNeutralContentRows = []models.ZoneContentRowSave{{Sid: "c", Count: 1}}
	state.HighNeutralContentRows = []models.ZoneContentRowSave{{Sid: "d", Count: 1}}
	state.HubZoneContentRows = []models.ZoneContentRowSave{{Sid: "e", Count: 1}}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	collectionLengths := []int{
		len(configuration.PlayerZoneMandatoryContent),
		len(configuration.LowNeutralMandatoryContent),
		len(configuration.MediumNeutralMandatoryContent),
		len(configuration.HighNeutralMandatoryContent),
		len(configuration.HubZoneMandatoryContent),
	}
	assert.Equal(t, []int{1, 1, 1, 1, 1}, collectionLengths)
}

func TestWhenPlayerRowHasCountTwo_ExpandsIntoTwoMandatoryItems(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "mine_gold", Count: 2, IsMine: true}}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Len(t, configuration.PlayerZoneMandatoryContent, 2)
}

func TestWhenPlayerRowIsMine_PropagatesIsMineFlag(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "mine_gold", Count: 1, IsMine: true}}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	require.Len(t, configuration.PlayerZoneMandatoryContent, 1)
	assert.True(t, configuration.PlayerZoneMandatoryContent[0].IsMine)
}

func TestWhenHighNeutralRowProvided_CopiesSidToMandatoryItem(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.HighNeutralContentRows = []models.ZoneContentRowSave{{Sid: "pandora_box", Count: 1}}

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	require.Len(t, configuration.HighNeutralMandatoryContent, 1)
	assert.Equal(t, "pandora_box", configuration.HighNeutralMandatoryContent[0].SID)
}

func TestWhenBonusEntriesProvided_CopiesBonuses(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	bonuses := []config_inner.BonusEntry{
		{PresetType: config_inner.BonusStartingWood, ReceiverFilter: "start_hero", Param: "7"},
	}
	state.BonusesJSON = bonuses

	// Act
	configuration := mappers.NewConfigMapper().FromEditorState(state)

	// Assert
	assert.Equal(t, bonuses, configuration.Bonuses)
}
