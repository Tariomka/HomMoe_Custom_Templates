package editorStateValidator_test

import (
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsDefault_ReturnsNoIssues(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()

	// Act
	issues := validate(&state)

	// Assert
	assert.Empty(t, issues)
}

func TestWhenStateHasInvalidValues_DoesNotModifyState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.PlayerCount = gofakeit.Number(9, 100)
	state.MapSize = 100
	state.NeutralZoneCount = gofakeit.Number(-100, -1)
	state.GameMode = "NotARealGameMode"
	state.VictoryCondition = "NotARealCondition"
	state.HeroCountMin = 10
	state.HeroCountMax = 3
	original := state

	// Act
	validate(&state)

	// Assert
	assert.Equal(t, original, state)
}

func TestWhenRangedFieldIsOutOfRange_ReturnsIssue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		mutate          func(state *editor_state_model.EditorState)
		expectedMessage string
	}{
		{"PlayerCountBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.PlayerCount = 1 },
			"playerCount 1 is outside [2, 8]"},
		{"PlayerCountAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.PlayerCount = 9 },
			"playerCount 9 is outside [2, 8]"},
		{"HeroMinBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.HeroCountMin = 0 },
			"heroMin 0 is outside [1, 12]"},
		{"HeroMaxAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.HeroCountMax = 13 },
			"heroMax 13 is outside [1, 12]"},
		{"HeroIncrementAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.HeroCountIncrement = 11 },
			"heroIncrement 11 is outside [1, 10]"},
		{"ResourceDensityBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.ResourceDensityPercent = 24 },
			"resourceDensity 24 is outside [25, 200]"},
		{"StructureDensityAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.StructureDensityPercent = 201 },
			"structureDensity 201 is outside [25, 200]"},
		{"NeutralStackStrengthAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.NeutralStackStrengthPercent = 300 },
			"neutralStackStrength 300 is outside [25, 200]"},
		{"BorderGuardStrengthBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.BorderGuardStrengthPercent = 0 },
			"borderGuardStrength 0 is outside [25, 200]"},
		{"FactionLawsExpBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.FactionLawXpPercent = 24 },
			"factionLawsExp 24 is outside [25, 200]"},
		{"AstrologyExpAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.AstrologyXpPercent = 500 },
			"astrologyExp 500 is outside [25, 200]"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			testCase.mutate(&state)

			// Act
			issues := validate(&state)

			// Assert
			assert.Contains(t, issueMessages(issues), testCase.expectedMessage)
		})
	}
}

func TestWhenCountFieldIsNegative_ReturnsIssue(t *testing.T) {
	t.Parallel()
	for _, testCase := range countFieldCases() {
		t.Run(testCase.name+"IsNegative_ReturnsIssue", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			negativeValue := gofakeit.Number(-1000, -1)
			testCase.mutate(&state, negativeValue)

			// Act
			issues := validate(&state)

			// Assert
			assert.Contains(t, issueMessages(issues),
				fmt.Sprintf("%s %d is outside [0, %d]", testCase.fieldName, negativeValue, testCase.highest))
		})
	}
}

func TestWhenCountFieldExceedsMaximum_ReturnsIssue(t *testing.T) {
	t.Parallel()
	for _, testCase := range countFieldCases() {
		t.Run(testCase.name+"ExceedsMaximum_ReturnsIssue", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			excessiveValue := gofakeit.Number(testCase.highest+1, testCase.highest+1000)
			testCase.mutate(&state, excessiveValue)

			// Act
			issues := validate(&state)

			// Assert
			assert.Contains(t, issueMessages(issues),
				fmt.Sprintf("%s %d is outside [0, %d]", testCase.fieldName, excessiveValue, testCase.highest))
		})
	}
}

func TestWhenCountFieldExceedsMaximum_FixClampsToMaximum(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.NeutralZoneCount = gofakeit.Number(17, 1000)

	// Act
	for _, issue := range validate(&state) {
		issue.Fix(&state)
	}

	// Assert
	assert.Equal(t, 16, state.NeutralZoneCount)
}

func TestWhenGameRuleFieldIsOutOfRange_ReturnsIssue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		mutate          func(state *editor_state_model.EditorState)
		expectedMessage string
	}{
		{"LostStartCityDayBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.LostStartCityDay = 0 },
			"lostStartCityDay 0 is outside [1, 30]"},
		{"CityHoldDaysAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.CityHoldDays = 31 },
			"cityHoldDays 31 is outside [1, 30]"},
		{"GladiatorArenaDaysDelayStartAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.GladiatorArenaDaysDelayStart = 61 },
			"gladiatorArenaDaysDelayStart 61 is outside [1, 60]"},
		{"GladiatorArenaCountDayBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.GladiatorArenaCountDay = 0 },
			"gladiatorArenaCountDay 0 is outside [1, 30]"},
		{"TournamentFirstTournamentDayBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.TournamentFirstTournamentDay = 2 },
			"tournamentFirstTournamentDay 2 is outside [3, 30]"},
		{"TournamentIntervalAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.TournamentInterval = 31 },
			"tournamentInterval 31 is outside [3, 30]"},
		{"TournamentPointsToWinAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.TournamentPointsToWin = 11 },
			"tournamentPointsToWin 11 is outside [1, 10]"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			testCase.mutate(&state)

			// Act
			issues := validate(&state)

			// Assert
			assert.Contains(t, issueMessages(issues), testCase.expectedMessage)
		})
	}
}

func TestWhenFloatFieldIsOutOfRange_ReturnsIssue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		mutate          func(state *editor_state_model.EditorState)
		expectedMessage string
	}{
		{"PlayerZoneSizeBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.PlayerZoneSize = 0.25 },
			"playerZoneSize 0.25 is outside [0.5, 2]"},
		{"NeutralZoneSizeAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.NeutralZoneSize = 2.5 },
			"neutralZoneSize 2.5 is outside [0.5, 2]"},
		{"HubZoneSizeBelowMinimum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.HubZoneSize = 0 },
			"hubZoneSize 0 is outside [0.5, 2]"},
		{"GuardRandomizationAboveMaximum_ReturnsIssue",
			func(state *editor_state_model.EditorState) { state.GuardRandomization = 0.75 },
			"guardRandomization 0.75 is outside [0, 0.5]"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			testCase.mutate(&state)

			// Act
			issues := validate(&state)

			// Assert
			assert.Contains(t, issueMessages(issues), testCase.expectedMessage)
		})
	}
}

func TestWhenFloatFieldIsOutOfRange_FixClampsToNearestBound(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.PlayerZoneSize = 0.25

	// Act
	for _, issue := range validate(&state) {
		issue.Fix(&state)
	}

	// Assert
	assert.InDelta(t, 0.5, state.PlayerZoneSize, 0.0001)
}

func TestWhenTopologyIsUnknown_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = "NotARealTopology"

	// Act
	issues := validate(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), `topology "NotARealTopology" is not a known topology`)
}

func TestWhenTopologyIsUnknown_FixRestoresRandom(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = "NotARealTopology"

	// Act
	for _, issue := range validate(&state) {
		issue.Fix(&state)
	}

	// Assert
	assert.Equal(t, config.TopologyRandom, state.Topology)
}

func TestWhenHeroMaxIsLessThanHeroMin_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.HeroCountMin = 6
	state.HeroCountMax = 5

	// Act
	issues := validate(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), "heroMax 5 is less than heroMin 6")
}

func TestWhenMapSizeIsUnknown_ReturnsIssueWithNearestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.MapSize = 100

	// Act
	issues := validate(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), "mapSize 100 is not a valid map size (nearest: 96)")
}

func TestWhenGameModeIsUnknown_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.GameMode = "NotARealGameMode"

	// Act
	issues := validate(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), `gameMode "NotARealGameMode" is not a known game mode`)
}

func TestWhenVictoryConditionIsUnknown_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.VictoryCondition = "NotARealCondition"

	// Act
	issues := validate(&state)

	// Assert
	assert.Contains(t, issueMessages(issues),
		`victoryCondition "NotARealCondition" is not a known victory condition`)
}

type countFieldCase struct {
	name      string
	fieldName string
	highest   int
	mutate    func(state *editor_state_model.EditorState, value int)
}

func countFieldCases() []countFieldCase {
	return []countFieldCase{
		{"NeutralZoneCount", "neutralZoneCount", 16,
			func(state *editor_state_model.EditorState, value int) { state.NeutralZoneCount = value }},
		{"PlayerOwnedCastles", "playerOwnedCastles", 4,
			func(state *editor_state_model.EditorState, value int) { state.PlayerOwnedCastles = value }},
		{"PlayerCastles", "playerCastles", 4,
			func(state *editor_state_model.EditorState, value int) { state.PlayerZoneCastles = value }},
		{"NeutralCastles", "neutralCastles", 4,
			func(state *editor_state_model.EditorState, value int) { state.NeutralZoneCastles = value }},
		{"AbandonedOutpostCount", "abandonedOutpostCount", 4,
			func(state *editor_state_model.EditorState, value int) { state.AbandonedOutpostCount = value }},
		{"HubCastles", "hubCastles", 4,
			func(state *editor_state_model.EditorState, value int) { state.HubZoneCastles = value }},
		{"RemoteFootholdCount", "remoteFootholdCount", 4,
			func(state *editor_state_model.EditorState, value int) { state.RemoteFootholdCount = value }},
		{"MaxPortalConns", "maxPortalConns", 32,
			func(state *editor_state_model.EditorState, value int) { state.MaxPortalConnections = value }},
		{"NeutralLowestNoCastle", "neutralLowestNoCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowestNoCastleCount = value }},
		{"NeutralLowestCastle", "neutralLowestCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowestCastleCount = value }},
		{"NeutralLowNoCastle", "neutralLowNoCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowNoCastleCount = value }},
		{"NeutralLowCastle", "neutralLowCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowCastleCount = value }},
		{"NeutralMediumNoCastle", "neutralMediumNoCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralMediumNoCastleCount = value }},
		{"NeutralMediumCastle", "neutralMediumCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralMediumCastleCount = value }},
		{"NeutralHighNoCastle", "neutralHighNoCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralHighNoCastleCount = value }},
		{"NeutralHighCastle", "neutralHighCastle", 8,
			func(state *editor_state_model.EditorState, value int) { state.NeutralHighCastleCount = value }},
		{"NeutralLowestCastlesPerZone", "neutralLowestCastlesPerZone", 4,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowestCastlesPerZone = value }},
		{"NeutralLowCastlesPerZone", "neutralLowCastlesPerZone", 4,
			func(state *editor_state_model.EditorState, value int) { state.NeutralLowCastlesPerZone = value }},
		{"NeutralMedCastlesPerZone", "neutralMedCastlesPerZone", 4,
			func(state *editor_state_model.EditorState, value int) { state.NeutralMediumCastlesPerZone = value }},
		{"NeutralHighCastlesPerZone", "neutralHighCastlesPerZone", 4,
			func(state *editor_state_model.EditorState, value int) { state.NeutralHighCastlesPerZone = value }},
	}
}

func issueMessages(issues []validators.ValidationIssue) []string {
	messages := make([]string, 0, len(issues))
	for _, issue := range issues {
		messages = append(messages, issue.Message)
	}
	return messages
}

func validate(state *editor_state_model.EditorState) []validators.ValidationIssue {
	return validators.NewEditorStateValidator().Validate(state)
}
