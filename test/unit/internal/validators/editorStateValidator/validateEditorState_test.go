package editorStateValidator_test

import (
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func issueMessages(issues []validators.ValidationIssue) []string {
	messages := make([]string, 0, len(issues))
	for _, issue := range issues {
		messages = append(messages, issue.Message)
	}
	return messages
}

func TestWhenStateIsDefault_ReturnsNoIssues(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()

	// Act
	issues := validators.ValidateEditorState(&state)

	// Assert
	assert.Empty(t, issues)
}

func TestWhenStateHasInvalidValues_DoesNotModifyState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerCount = gofakeit.Number(9, 100)
	state.MapSize = 100
	state.NeutralZoneCount = gofakeit.Number(-100, -1)
	state.GameMode = "NotARealGameMode"
	state.VictoryCondition = "NotARealCondition"
	state.HeroCountMin = 10
	state.HeroCountMax = 3
	original := state

	// Act
	validators.ValidateEditorState(&state)

	// Assert
	assert.Equal(t, original, state)
}

func TestWhenRangedFieldIsOutOfRange_ReturnsIssue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		mutate          func(state *dtos.EditorStateDto)
		expectedMessage string
	}{
		{"PlayerCountBelowMinimum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.PlayerCount = 1 },
			"playerCount 1 is outside [2, 8]"},
		{"PlayerCountAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.PlayerCount = 9 },
			"playerCount 9 is outside [2, 8]"},
		{"HeroMinBelowMinimum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.HeroCountMin = 0 },
			"heroMin 0 is outside [1, 12]"},
		{"HeroMaxAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.HeroCountMax = 13 },
			"heroMax 13 is outside [1, 12]"},
		{"HeroIncrementAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.HeroCountIncrement = 11 },
			"heroIncrement 11 is outside [1, 10]"},
		{"ResourceDensityBelowMinimum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.ResourceDensityPercent = 24 },
			"resourceDensity 24 is outside [25, 200]"},
		{"StructureDensityAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.StructureDensityPercent = 201 },
			"structureDensity 201 is outside [25, 200]"},
		{"NeutralStackStrengthAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.NeutralStackStrengthPercent = 300 },
			"neutralStackStrength 300 is outside [25, 200]"},
		{"BorderGuardStrengthBelowMinimum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.BorderGuardStrengthPercent = 0 },
			"borderGuardStrength 0 is outside [25, 200]"},
		{"FactionLawsExpBelowMinimum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.FactionLawXpPercent = 24 },
			"factionLawsExp 24 is outside [25, 200]"},
		{"AstrologyExpAboveMaximum_ReturnsIssue",
			func(state *dtos.EditorStateDto) { state.AstrologyXpPercent = 500 },
			"astrologyExp 500 is outside [25, 200]"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := dtos.NewDefaultEditorStateDto()
			testCase.mutate(&state)

			// Act
			issues := validators.ValidateEditorState(&state)

			// Assert
			assert.Contains(t, issueMessages(issues), testCase.expectedMessage)
		})
	}
}

func TestWhenCountFieldIsNegative_ReturnsIssue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		fieldName string
		mutate    func(state *dtos.EditorStateDto, value int)
	}{
		{"NeutralZoneCountIsNegative_ReturnsIssue", "neutralZoneCount",
			func(state *dtos.EditorStateDto, value int) { state.NeutralZoneCount = value }},
		{"PlayerOwnedCastlesIsNegative_ReturnsIssue", "playerOwnedCastles",
			func(state *dtos.EditorStateDto, value int) { state.PlayerOwnedCastles = value }},
		{"PlayerCastlesIsNegative_ReturnsIssue", "playerCastles",
			func(state *dtos.EditorStateDto, value int) { state.PlayerZoneCastles = value }},
		{"NeutralCastlesIsNegative_ReturnsIssue", "neutralCastles",
			func(state *dtos.EditorStateDto, value int) { state.NeutralZoneCastles = value }},
		{"AbandonedOutpostCountIsNegative_ReturnsIssue", "abandonedOutpostCount",
			func(state *dtos.EditorStateDto, value int) { state.AbandonedOutpostCount = value }},
		{"NeutralLowNoCastleIsNegative_ReturnsIssue", "neutralLowNoCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralLowNoCastleCount = value }},
		{"NeutralLowCastleIsNegative_ReturnsIssue", "neutralLowCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralLowCastleCount = value }},
		{"NeutralMediumNoCastleIsNegative_ReturnsIssue", "neutralMediumNoCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralMediumNoCastleCount = value }},
		{"NeutralMediumCastleIsNegative_ReturnsIssue", "neutralMediumCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralMediumCastleCount = value }},
		{"NeutralHighNoCastleIsNegative_ReturnsIssue", "neutralHighNoCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralHighNoCastleCount = value }},
		{"NeutralHighCastleIsNegative_ReturnsIssue", "neutralHighCastle",
			func(state *dtos.EditorStateDto, value int) { state.NeutralHighCastleCount = value }},
		{"NeutralLowCastlesPerZoneIsNegative_ReturnsIssue", "neutralLowCastlesPerZone",
			func(state *dtos.EditorStateDto, value int) { state.NeutralLowCastlesPerZone = value }},
		{"NeutralMedCastlesPerZoneIsNegative_ReturnsIssue", "neutralMedCastlesPerZone",
			func(state *dtos.EditorStateDto, value int) { state.NeutralMediumCastlesPerZone = value }},
		{"NeutralHighCastlesPerZoneIsNegative_ReturnsIssue", "neutralHighCastlesPerZone",
			func(state *dtos.EditorStateDto, value int) { state.NeutralHighCastlesPerZone = value }},
		{"HubCastlesIsNegative_ReturnsIssue", "hubCastles",
			func(state *dtos.EditorStateDto, value int) { state.HubZoneCastles = value }},
		{"MinNeutralZonesBetweenPlayersIsNegative_ReturnsIssue", "minNeutralZonesBetweenPlayers",
			func(state *dtos.EditorStateDto, value int) { state.MinNeutralZonesBetweenPlayers = value }},
		{"RemoteFootholdCountIsNegative_ReturnsIssue", "remoteFootholdCount",
			func(state *dtos.EditorStateDto, value int) { state.RemoteFootholdCount = value }},
		{"MaxPortalConnsIsNegative_ReturnsIssue", "maxPortalConns",
			func(state *dtos.EditorStateDto, value int) { state.MaxPortalConnections = value }},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := dtos.NewDefaultEditorStateDto()
			negativeValue := gofakeit.Number(-1000, -1)
			testCase.mutate(&state, negativeValue)

			// Act
			issues := validators.ValidateEditorState(&state)

			// Assert
			assert.Contains(t, issueMessages(issues),
				fmt.Sprintf("%s %d is negative", testCase.fieldName, negativeValue))
		})
	}
}

func TestWhenHeroMaxIsLessThanHeroMin_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.HeroCountMin = 6
	state.HeroCountMax = 5

	// Act
	issues := validators.ValidateEditorState(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), "heroMax 5 is less than heroMin 6")
}

func TestWhenMapSizeIsUnknown_ReturnsIssueWithNearestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.MapSize = 100

	// Act
	issues := validators.ValidateEditorState(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), "mapSize 100 is not a valid map size (nearest: 96)")
}

func TestWhenGameModeIsUnknown_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.GameMode = "NotARealGameMode"

	// Act
	issues := validators.ValidateEditorState(&state)

	// Assert
	assert.Contains(t, issueMessages(issues), `gameMode "NotARealGameMode" is not a known game mode`)
}

func TestWhenVictoryConditionIsUnknown_ReturnsIssue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.VictoryCondition = "NotARealCondition"

	// Act
	issues := validators.ValidateEditorState(&state)

	// Assert
	assert.Contains(t, issueMessages(issues),
		`victoryCondition "NotARealCondition" is not a known victory condition`)
}
