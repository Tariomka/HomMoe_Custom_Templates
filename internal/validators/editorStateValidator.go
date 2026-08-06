package validators

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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
func (this *EditorStateValidator) Validate(state *dtos.EditorStateDto) []ValidationIssue {
	issues := this.validateRangedFields(state)
	issues = append(issues, this.validateRangedFloatFields(state)...)
	issues = append(issues, this.validateHeroCountOrder(state)...)
	issues = append(issues, this.validateMapSize(state)...)
	issues = append(issues, this.validateGameMode(state)...)
	issues = append(issues, this.validateVictoryCondition(state)...)
	issues = append(issues, this.validateTopology(state)...)
	return issues
}

func (this *EditorStateValidator) validateRangedFields(state *dtos.EditorStateDto) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, check := range this.rangedIntFields() {
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

func (this *EditorStateValidator) validateRangedFloatFields(state *dtos.EditorStateDto) []ValidationIssue {
	issues := []ValidationIssue{}
	for _, check := range this.rangedFloatFields() {
		currentValue := *check.field.value(state)
		if currentValue >= check.lowest && currentValue <= check.highest {
			continue
		}
		issues = append(issues, ValidationIssue{
			Message: fmt.Sprintf(
				"%s %g is outside [%g, %g]", check.field.name, currentValue, check.lowest, check.highest),
			fix: func(state *dtos.EditorStateDto) {
				fieldValue := check.field.value(state)
				*fieldValue = helpers.Clamp(*fieldValue, check.lowest, check.highest)
			},
		})
	}
	return issues
}

func (this *EditorStateValidator) validateHeroCountOrder(state *dtos.EditorStateDto) []ValidationIssue {
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

func (this *EditorStateValidator) validateMapSize(state *dtos.EditorStateDto) []ValidationIssue {
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

func (this *EditorStateValidator) validateGameMode(state *dtos.EditorStateDto) []ValidationIssue {
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

func (this *EditorStateValidator) validateVictoryCondition(state *dtos.EditorStateDto) []ValidationIssue {
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

func (this *EditorStateValidator) validateTopology(state *dtos.EditorStateDto) []ValidationIssue {
	for descriptor := range common_topologies.GetTopologyDescriptorSeq() {
		if descriptor.Type == state.Topology {
			return nil
		}
	}
	return []ValidationIssue{{
		Message: fmt.Sprintf("topology %q is not a known topology", state.Topology),
		fix: func(state *dtos.EditorStateDto) {
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
func (this *EditorStateValidator) rangedIntFields() []rangedIntField {
	return []rangedIntField{
		newRangedIntField("playerCount", playerCountLowest, playerCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.PlayerCount }),
		newRangedIntField("heroMin", heroCountLowest, heroCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.HeroCountMin }),
		newRangedIntField("heroMax", heroCountLowest, heroCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.HeroCountMax }),
		newRangedIntField("heroIncrement", heroIncrementLowest, heroIncrementHighest,
			func(state *dtos.EditorStateDto) *int { return &state.HeroCountIncrement }),
		newRangedIntField("resourceDensity", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.ResourceDensityPercent }),
		newRangedIntField("structureDensity", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.StructureDensityPercent }),
		newRangedIntField("neutralStackStrength", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralStackStrengthPercent }),
		newRangedIntField("borderGuardStrength", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.BorderGuardStrengthPercent }),
		newRangedIntField("factionLawsExp", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.FactionLawXpPercent }),
		newRangedIntField("astrologyExp", percentLowest, percentHighest,
			func(state *dtos.EditorStateDto) *int { return &state.AstrologyXpPercent }),
		newRangedIntField("neutralZoneCount", countLowest, neutralZoneCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralZoneCount }),
		newRangedIntField("playerOwnedCastles", countLowest, castleCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.PlayerOwnedCastles }),
		newRangedIntField("playerCastles", countLowest, castleCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.PlayerZoneCastles }),
		newRangedIntField("neutralCastles", countLowest, castleCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralZoneCastles }),
		newRangedIntField("abandonedOutpostCount", countLowest, castleCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.AbandonedOutpostCount }),
		newRangedIntField("hubCastles", countLowest, castleCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.HubZoneCastles }),
		newRangedIntField("remoteFootholdCount", countLowest, remoteFootholdHighest,
			func(state *dtos.EditorStateDto) *int { return &state.RemoteFootholdCount }),
		newRangedIntField("maxPortalConns", countLowest, portalConnectionsHighest,
			func(state *dtos.EditorStateDto) *int { return &state.MaxPortalConnections }),
		newRangedIntField("neutralLowestNoCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestNoCastleCount }),
		newRangedIntField("neutralLowestCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestCastleCount }),
		newRangedIntField("neutralLowNoCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowNoCastleCount }),
		newRangedIntField("neutralLowCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowCastleCount }),
		newRangedIntField("neutralMediumNoCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumNoCastleCount }),
		newRangedIntField("neutralMediumCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumCastleCount }),
		newRangedIntField("neutralHighNoCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralHighNoCastleCount }),
		newRangedIntField("neutralHighCastle", countLowest, neutralTierCountHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralHighCastleCount }),
		newRangedIntField("neutralLowestCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowestCastlesPerZone }),
		newRangedIntField("neutralLowCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralLowCastlesPerZone }),
		newRangedIntField("neutralMedCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralMediumCastlesPerZone }),
		newRangedIntField("neutralHighCastlesPerZone", countLowest, castlesPerZoneHighest,
			func(state *dtos.EditorStateDto) *int { return &state.NeutralHighCastlesPerZone }),
		newRangedIntField("lostStartCityDay", ruleDayLowest, ruleDayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.LostStartCityDay }),
		newRangedIntField("cityHoldDays", ruleDayLowest, ruleDayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.CityHoldDays }),
		newRangedIntField("gladiatorArenaDaysDelayStart", ruleDayLowest, arenaDelayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.GladiatorArenaDaysDelayStart }),
		newRangedIntField("gladiatorArenaCountDay", ruleDayLowest, ruleDayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.GladiatorArenaCountDay }),
		newRangedIntField("tournamentFirstTournamentDay", tournamentDayLowest, tournamentDayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.TournamentFirstTournamentDay }),
		newRangedIntField("tournamentInterval", tournamentDayLowest, tournamentDayHighest,
			func(state *dtos.EditorStateDto) *int { return &state.TournamentInterval }),
		newRangedIntField("tournamentPointsToWin", tournamentPointsLowest, tournamentPointsHighest,
			func(state *dtos.EditorStateDto) *int { return &state.TournamentPointsToWin }),
	}
}

func (this *EditorStateValidator) rangedFloatFields() []rangedFloatField {
	return []rangedFloatField{
		newRangedFloatField("playerZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *dtos.EditorStateDto) *float64 { return &state.PlayerZoneSize }),
		newRangedFloatField("neutralZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *dtos.EditorStateDto) *float64 { return &state.NeutralZoneSize }),
		newRangedFloatField("hubZoneSize", zoneSizeLowest, zoneSizeHighest,
			func(state *dtos.EditorStateDto) *float64 { return &state.HubZoneSize }),
		newRangedFloatField("guardRandomization", guardRandomizationLowest, guardRandomizationHighest,
			func(state *dtos.EditorStateDto) *float64 { return &state.GuardRandomization }),
	}
}
