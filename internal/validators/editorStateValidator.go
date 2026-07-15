package validators

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const (
	playerCountLowest    = 2
	playerCountHighest   = 8
	heroCountLowest      = 1
	heroCountHighest     = 12
	heroIncrementLowest  = 1
	heroIncrementHighest = 10
	percentLowest        = 25
	percentHighest       = 200
)

// intField names a .gen.json integer field and provides access to it, so a
// single check can both read the value (validation) and write it (fix).
type intField struct {
	name  string
	value func(state *dtos.EditorStateDto) *int
}

// ValidateEditorState checks an editor state loaded from a .gen.json file
// against the editor's allowed ranges and the known registry values. It never
// modifies the state; each returned issue carries a Fix that applies the
// correction. An empty result means the state is valid.
func ValidateEditorState(state *dtos.EditorStateDto) []ValidationIssue {
	issues := validateRangedFields(state)
	issues = append(issues, validateNonNegativeFields(state)...)
	issues = append(issues, validateHeroCountOrder(state)...)
	issues = append(issues, validateMapSize(state)...)
	issues = append(issues, validateGameMode(state)...)
	issues = append(issues, validateVictoryCondition(state)...)
	return issues
}

func validateRangedFields(state *dtos.EditorStateDto) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, check := range rangedIntFields() {
		currentValue := *check.field.value(state)
		if currentValue >= check.lowest && currentValue <= check.highest {
			continue
		}
		issues = append(issues, ValidationIssue{
			Message: fmt.Sprintf(
				"%s %d is outside [%d, %d]", check.field.name, currentValue, check.lowest, check.highest),
			fix: func(state *dtos.EditorStateDto) {
				fieldValue := check.field.value(state)
				*fieldValue = helpers.Clamp(*fieldValue, check.lowest, check.highest)
			},
		})
	}
	return issues
}

func validateNonNegativeFields(state *dtos.EditorStateDto) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, field := range nonNegativeIntFields() {
		currentValue := *field.value(state)
		if currentValue >= 0 {
			continue
		}
		issues = append(issues, ValidationIssue{
			Message: fmt.Sprintf("%s %d is negative", field.name, currentValue),
			fix: func(state *dtos.EditorStateDto) {
				fieldValue := field.value(state)
				*fieldValue = max(*fieldValue, 0)
			},
		})
	}
	return issues
}

func validateHeroCountOrder(state *dtos.EditorStateDto) []ValidationIssue {
	if state.HeroCountMax >= state.HeroCountMin {
		return nil
	}
	return []ValidationIssue{{
		Message: fmt.Sprintf("heroMax %d is less than heroMin %d", state.HeroCountMax, state.HeroCountMin),
		fix: func(state *dtos.EditorStateDto) {
			state.HeroCountMax = max(state.HeroCountMax, state.HeroCountMin)
		},
	}}
}

func validateMapSize(state *dtos.EditorStateDto) []ValidationIssue {
	nearest := common.GetNearestMapSize(state.MapSize)
	if nearest.Size == state.MapSize {
		return nil
	}
	return []ValidationIssue{{
		Message: fmt.Sprintf("mapSize %d is not a valid map size (nearest: %d)", state.MapSize, nearest.Size),
		fix: func(state *dtos.EditorStateDto) {
			state.MapSize = common.GetNearestMapSize(state.MapSize).Size
		},
	}}
}

func validateGameMode(state *dtos.EditorStateDto) []ValidationIssue {
	if slices.Contains(registry.GetGameModeList(), state.GameMode) {
		return nil
	}
	return []ValidationIssue{{
		Message: fmt.Sprintf("gameMode %q is not a known game mode", state.GameMode),
		fix: func(state *dtos.EditorStateDto) {
			state.GameMode = registry.GetGameModeValues().Classic
		},
	}}
}

func validateVictoryCondition(state *dtos.EditorStateDto) []ValidationIssue {
	winConditionValues := registry.GetWinningConditionValues()
	validConditions := []string{
		winConditionValues.Standard,
		winConditionValues.CapitalCapture,
		winConditionValues.CapitalHold,
		winConditionValues.FinalBattle,
		winConditionValues.CityHold,
		winConditionValues.Tournament,
	}
	if slices.Contains(validConditions, state.VictoryCondition) {
		return nil
	}
	return []ValidationIssue{{
		Message: fmt.Sprintf("victoryCondition %q is not a known victory condition", state.VictoryCondition),
		fix: func(state *dtos.EditorStateDto) {
			state.VictoryCondition = registry.GetWinningConditionValues().Standard
		},
	}}
}

type rangedIntField struct {
	field   intField
	lowest  int
	highest int
}

func rangedIntFields() []rangedIntField {
	return []rangedIntField{
		{intField{"playerCount", func(state *dtos.EditorStateDto) *int { return &state.PlayerCount }},
			playerCountLowest, playerCountHighest},
		{intField{"heroMin", func(state *dtos.EditorStateDto) *int { return &state.HeroCountMin }},
			heroCountLowest, heroCountHighest},
		{intField{"heroMax", func(state *dtos.EditorStateDto) *int { return &state.HeroCountMax }},
			heroCountLowest, heroCountHighest},
		{intField{"heroIncrement", func(state *dtos.EditorStateDto) *int { return &state.HeroCountIncrement }},
			heroIncrementLowest, heroIncrementHighest},
		{intField{"resourceDensity", func(state *dtos.EditorStateDto) *int { return &state.ResourceDensityPercent }},
			percentLowest, percentHighest},
		{intField{"structureDensity", func(state *dtos.EditorStateDto) *int { return &state.StructureDensityPercent }},
			percentLowest, percentHighest},
		{
			intField{
				"neutralStackStrength",
				func(state *dtos.EditorStateDto) *int { return &state.NeutralStackStrengthPercent },
			},
			percentLowest,
			percentHighest,
		},
		{
			intField{
				"borderGuardStrength",
				func(state *dtos.EditorStateDto) *int { return &state.BorderGuardStrengthPercent },
			},
			percentLowest,
			percentHighest,
		},
		{intField{"factionLawsExp", func(state *dtos.EditorStateDto) *int { return &state.FactionLawXpPercent }},
			percentLowest, percentHighest},
		{intField{"astrologyExp", func(state *dtos.EditorStateDto) *int { return &state.AstrologyXpPercent }},
			percentLowest, percentHighest},
	}
}

func nonNegativeIntFields() []intField {
	return []intField{
		{"neutralZoneCount", func(state *dtos.EditorStateDto) *int { return &state.NeutralZoneCount }},
		{"playerOwnedCastles", func(state *dtos.EditorStateDto) *int { return &state.PlayerOwnedCastles }},
		{"playerCastles", func(state *dtos.EditorStateDto) *int { return &state.PlayerZoneCastles }},
		{"neutralCastles", func(state *dtos.EditorStateDto) *int { return &state.NeutralZoneCastles }},
		{"abandonedOutpostCount", func(state *dtos.EditorStateDto) *int { return &state.AbandonedOutpostCount }},
		{"neutralLowestNoCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestNoCastleCount }},
		{"neutralLowestCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestCastleCount }},
		{"neutralLowNoCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralLowNoCastleCount }},
		{"neutralLowCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralLowCastleCount }},
		{"neutralMediumNoCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumNoCastleCount }},
		{"neutralMediumCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumCastleCount }},
		{"neutralHighNoCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralHighNoCastleCount }},
		{"neutralHighCastle", func(state *dtos.EditorStateDto) *int { return &state.NeutralHighCastleCount }},
		{
			"neutralLowestCastlesPerZone",
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestCastlesPerZone },
		},
		{"neutralLowCastlesPerZone", func(state *dtos.EditorStateDto) *int { return &state.NeutralLowCastlesPerZone }},
		{
			"neutralMedCastlesPerZone",
			func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumCastlesPerZone },
		},
		{
			"neutralHighCastlesPerZone",
			func(state *dtos.EditorStateDto) *int { return &state.NeutralHighCastlesPerZone },
		},
		{"hubCastles", func(state *dtos.EditorStateDto) *int { return &state.HubZoneCastles }},
		{"remoteFootholdCount", func(state *dtos.EditorStateDto) *int { return &state.RemoteFootholdCount }},
		{"maxPortalConns", func(state *dtos.EditorStateDto) *int { return &state.MaxPortalConnections }},
	}
}
