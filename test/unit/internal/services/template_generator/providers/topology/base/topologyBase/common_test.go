package topologyBase_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// newUnitTuning builds tuning with every multiplier at 1.0 so that guard and
// content values pass through the scaling helpers unchanged, keeping the
// expected values in assertions readable.
func newUnitTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      1.0,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
		GuardRandomization:             0.05,
	}
}

// newSpawnRequest builds a spawn-zone request with the invariant fields the
// tests never vary (unit zone size, no faction matching, Player1) already
// filled in.
func newSpawnRequest(
	label string,
	connectionNames []string,
	castleCount, footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning) models.SpawnZoneCreationRequest {
	return models.SpawnZoneCreationRequest{
		Label:           label,
		PlayerName:      "Player1",
		ConnectionNames: connectionNames,
		CastleCount:     castleCount,
		MatchFactions:   false,
		Size:            1.0,
		FootholdCount:   footholdCount,
		GenerateRoads:   generateRoads,
		Tuning:          tuning,
	}
}

// newNeutralRequest builds a neutral-zone request with the invariant fields the
// tests never vary (unit zone size, roads enabled) already filled in.
func newNeutralRequest(
	plan neutral_zone.Plan,
	connectionNames []string,
	footholdCount int,
	tuning models.GenerationTuning,
	holdCity bool) models.TopologyNeutralZoneCreationRequest {
	return models.TopologyNeutralZoneCreationRequest{
		Plan:            plan,
		ConnectionNames: connectionNames,
		Size:            1.0,
		FootholdCount:   footholdCount,
		GenerateRoads:   true,
		HoldCity:        holdCity,
		Tuning:          tuning,
	}
}
