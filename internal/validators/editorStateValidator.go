package validators

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
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

	// Count ceilings mirror the editor's sliders. Their floor stays at zero
	// rather than at the slider minimum so files that were valid before these
	// ceilings existed keep loading without new warnings.
	countLowest              = 0
	neutralZoneCountHighest  = 16
	castleCountHighest       = 4
	neutralTierCountHighest  = 8
	castlesPerZoneHighest    = 4
	remoteFootholdHighest    = 4
	portalConnectionsHighest = 32

	ruleDayLowest           = 1
	ruleDayHighest          = 30
	arenaDelayHighest       = 60
	tournamentDayLowest     = 3
	tournamentDayHighest    = 30
	tournamentPointsLowest  = 1
	tournamentPointsHighest = 10

	zoneSizeLowest            = 0.5
	zoneSizeHighest           = 2.0
	guardRandomizationLowest  = 0.0
	guardRandomizationHighest = 0.5
)

type EditorStateValidator struct{}

func NewEditorStateValidator() IEditorStateValidator {
	return &EditorStateValidator{}
}

// Validate checks an editor state loaded from a .gen.json file
// against the editor's allowed ranges and the known registry values. It never
// modifies the state; each returned issue carries a Fix that applies the
// correction. An empty result means the state is valid.
func (this *EditorStateValidator) Validate(state *editor_state_model.EditorState) []ValidationIssue {
	issues := this.validateRangedFields(state)
	issues = append(issues, this.validateRangedFloatFields(state)...)
	issues = append(issues, this.validateHeroCountOrder(state)...)
	issues = append(issues, this.validateMapSize(state)...)
	issues = append(issues, this.validateGameMode(state)...)
	issues = append(issues, this.validateVictoryCondition(state)...)
	issues = append(issues, this.validateTopology(state)...)
	return issues
}

func (this *EditorStateValidator) validateRangedFields(state *editor_state_model.EditorState) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, check := range this.getRangedIntFieldDescriptors() {
		currentValue := *check.value(state)
		if currentValue >= check.lowest && currentValue <= check.highest {
			continue
		}

		issues = append(issues, ValidationIssue{
			Message: fmt.Sprintf("%s %d is outside [%d, %d]", check.name, currentValue, check.lowest, check.highest),
			fix: func(state *editor_state_model.EditorState) {
				fieldValue := check.value(state)
				*fieldValue = helpers.Clamp(*fieldValue, check.lowest, check.highest)
			},
		})
	}
	return issues
}

func (this *EditorStateValidator) validateRangedFloatFields(
	state *editor_state_model.EditorState) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, check := range this.getRangedFloatFieldDescriptors() {
		currentValue := *check.value(state)
		if currentValue >= check.lowest && currentValue <= check.highest {
			continue
		}

		issues = append(issues, ValidationIssue{
			Message: fmt.Sprintf("%s %g is outside [%g, %g]", check.name, currentValue, check.lowest, check.highest),
			fix: func(state *editor_state_model.EditorState) {
				fieldValue := check.value(state)
				*fieldValue = helpers.Clamp(*fieldValue, check.lowest, check.highest)
			},
		})
	}
	return issues
}

func (this *EditorStateValidator) validateHeroCountOrder(state *editor_state_model.EditorState) []ValidationIssue {
	if state.HeroCountMax >= state.HeroCountMin {
		return nil
	}

	return []ValidationIssue{{
		Message: fmt.Sprintf("heroMax %d is less than heroMin %d", state.HeroCountMax, state.HeroCountMin),
		fix: func(state *editor_state_model.EditorState) {
			state.HeroCountMax = max(state.HeroCountMax, state.HeroCountMin)
		},
	}}
}

func (this *EditorStateValidator) validateMapSize(state *editor_state_model.EditorState) []ValidationIssue {
	nearest := common.GetNearestMapSize(state.MapSize)
	if nearest.Size == state.MapSize {
		return nil
	}

	return []ValidationIssue{{
		Message: fmt.Sprintf("mapSize %d is not a valid map size (nearest: %d)", state.MapSize, nearest.Size),
		fix: func(state *editor_state_model.EditorState) {
			state.MapSize = common.GetNearestMapSize(state.MapSize).Size
		},
	}}
}

func (this *EditorStateValidator) validateGameMode(state *editor_state_model.EditorState) []ValidationIssue {
	if slices.Contains(registry.GetGameModeList(), state.GameMode) {
		return nil
	}

	return []ValidationIssue{{
		Message: fmt.Sprintf("gameMode %q is not a known game mode", state.GameMode),
		fix: func(state *editor_state_model.EditorState) {
			state.GameMode = registry.GetGameModeValues().Classic
		},
	}}
}

func (this *EditorStateValidator) validateVictoryCondition(
	state *editor_state_model.EditorState) []ValidationIssue {
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
		fix: func(state *editor_state_model.EditorState) {
			state.VictoryCondition = registry.GetWinningConditionValues().Standard
		},
	}}
}

func (this *EditorStateValidator) validateTopology(state *editor_state_model.EditorState) []ValidationIssue {
	for descriptor := range common_topologies.GetTopologyDescriptorSeq() {
		if descriptor.Type == state.Topology {
			return nil
		}
	}

	return []ValidationIssue{{
		Message: fmt.Sprintf("topology %q is not a known topology", state.Topology),
		fix: func(state *editor_state_model.EditorState) {
			state.Topology = config.TopologyRandom
		},
	}}
}

// rangedIntFields lists every integer field of the editor state together with
// the bounds its editor control offers. Count fields keep a zero floor rather
// than their control's minimum so states written before these ceilings existed
// keep loading unchanged.
//
//nolint:funlen // Declarative field table.
func (this *EditorStateValidator) getRangedIntFieldDescriptors() []rangedFieldDescriptor[int] {
	return []rangedFieldDescriptor[int]{
		newRangedFieldDescriptor("playerCount", playerCountLowest, playerCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.PlayerCount }),
		newRangedFieldDescriptor("heroMin", heroCountLowest, heroCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.HeroCountMin }),
		newRangedFieldDescriptor("heroMax", heroCountLowest, heroCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.HeroCountMax }),
		newRangedFieldDescriptor("heroIncrement", heroIncrementLowest, heroIncrementHighest,
			func(state *editor_state_model.EditorState) *int { return &state.HeroCountIncrement }),
		newRangedFieldDescriptor("resourceDensity", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.ResourceDensityPercent }),
		newRangedFieldDescriptor("structureDensity", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.StructureDensityPercent }),
		newRangedFieldDescriptor("neutralStackStrength", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralStackStrengthPercent }),
		newRangedFieldDescriptor("borderGuardStrength", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.BorderGuardStrengthPercent }),
		newRangedFieldDescriptor("factionLawsExp", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.FactionLawXpPercent }),
		newRangedFieldDescriptor("astrologyExp", percentLowest, percentHighest,
			func(state *editor_state_model.EditorState) *int { return &state.AstrologyXpPercent }),
		newRangedFieldDescriptor("neutralZoneCount", countLowest, neutralZoneCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralZoneCount }),
		newRangedFieldDescriptor("playerOwnedCastles", countLowest, castleCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.PlayerOwnedCastles }),
		newRangedFieldDescriptor("playerCastles", countLowest, castleCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.PlayerZoneCastles }),
		newRangedFieldDescriptor("neutralCastles", countLowest, castleCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralZoneCastles }),
		newRangedFieldDescriptor("abandonedOutpostCount", countLowest, castleCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.AbandonedOutpostCount }),
		newRangedFieldDescriptor("hubCastles", countLowest, castleCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.HubZoneCastles }),
		newRangedFieldDescriptor("remoteFootholdCount", countLowest, remoteFootholdHighest,
			func(state *editor_state_model.EditorState) *int { return &state.RemoteFootholdCount }),
		newRangedFieldDescriptor("maxPortalConns", countLowest, portalConnectionsHighest,
			func(state *editor_state_model.EditorState) *int { return &state.MaxPortalConnections }),
		newRangedFieldDescriptor("neutralLowestNoCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowestNoCastleCount }),
		newRangedFieldDescriptor("neutralLowestCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowestCastleCount }),
		newRangedFieldDescriptor("neutralLowNoCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowNoCastleCount }),
		newRangedFieldDescriptor("neutralLowCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowCastleCount }),
		newRangedFieldDescriptor("neutralMediumNoCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralMediumNoCastleCount }),
		newRangedFieldDescriptor("neutralMediumCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralMediumCastleCount }),
		newRangedFieldDescriptor("neutralHighNoCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralHighNoCastleCount }),
		newRangedFieldDescriptor("neutralHighCastle", countLowest, neutralTierCountHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralHighCastleCount }),
		newRangedFieldDescriptor("neutralLowestCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowestCastlesPerZone }),
		newRangedFieldDescriptor("neutralLowCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralLowCastlesPerZone }),
		newRangedFieldDescriptor("neutralMedCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralMediumCastlesPerZone }),
		newRangedFieldDescriptor("neutralHighCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *editor_state_model.EditorState) *int { return &state.NeutralHighCastlesPerZone }),
		newRangedFieldDescriptor("lostStartCityDay", ruleDayLowest, ruleDayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.LostStartCityDay }),
		newRangedFieldDescriptor("cityHoldDays", ruleDayLowest, ruleDayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.CityHoldDays }),
		newRangedFieldDescriptor("gladiatorArenaDaysDelayStart", ruleDayLowest, arenaDelayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.GladiatorArenaDaysDelayStart }),
		newRangedFieldDescriptor("gladiatorArenaCountDay", ruleDayLowest, ruleDayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.GladiatorArenaCountDay }),
		newRangedFieldDescriptor("tournamentFirstTournamentDay", tournamentDayLowest, tournamentDayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.TournamentFirstTournamentDay }),
		newRangedFieldDescriptor("tournamentInterval", tournamentDayLowest, tournamentDayHighest,
			func(state *editor_state_model.EditorState) *int { return &state.TournamentInterval }),
		newRangedFieldDescriptor("tournamentPointsToWin", tournamentPointsLowest, tournamentPointsHighest,
			func(state *editor_state_model.EditorState) *int { return &state.TournamentPointsToWin }),
	}
}

func (this *EditorStateValidator) getRangedFloatFieldDescriptors() []rangedFieldDescriptor[float64] {
	return []rangedFieldDescriptor[float64]{
		newRangedFieldDescriptor("playerZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *editor_state_model.EditorState) *float64 { return &state.PlayerZoneSize }),
		newRangedFieldDescriptor("neutralZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *editor_state_model.EditorState) *float64 { return &state.NeutralZoneSize }),
		newRangedFieldDescriptor("hubZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *editor_state_model.EditorState) *float64 { return &state.HubZoneSize }),
		newRangedFieldDescriptor("guardRandomization", guardRandomizationLowest, guardRandomizationHighest,
			func(state *editor_state_model.EditorState) *float64 { return &state.GuardRandomization }),
	}
}
