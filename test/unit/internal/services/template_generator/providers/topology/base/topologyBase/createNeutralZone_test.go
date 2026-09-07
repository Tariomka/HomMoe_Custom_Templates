package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsCreated_NameCombinesNeutralPrefixWithPlanLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(3)
	plan := neutral_zone.Plan{Label: label, Quality: neutral_zone.QualityMedium, CastleCount: 1}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), false))

	// Assert
	assert.Equal(t, "Neutral-"+label, zone.Name)
}

func TestWhenPlanHasCastles_MainObjectCountMatchesCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 2}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), false))

	// Assert
	assert.Len(t, zone.MainObjects, 2)
}

func TestWhenHoldCityZoneHasNoPlannedCastles_ForcesSingleCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 0}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), true))

	// Assert
	assert.Len(t, zone.MainObjects, 1)
}

func TestWhenZoneIsHoldCity_PrimaryCastleCarriesHoldCityWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityHigh, CastleCount: 1}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), true))

	// Assert
	assert.True(t, zone.MainObjects[0].HoldCityWinCon)
}

func TestWhenZoneHasNoMainObjects_BiomeMatchesZone(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 0}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), false))

	// Assert
	assert.Equal(t, template_model.TypedRef{Type: "MatchZone"}, zone.ZoneBiome)
}

func TestWhenZoneHasMainObjects_BiomeMatchesPrimaryMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 1}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, newUnitTuning(), false))

	// Assert
	assert.Equal(t, template_model.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}, zone.ZoneBiome)
}

func TestWhenAbandonedOutpostCountIsPositive_AppendsOutpostsAfterCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 1}
	tuning := newUnitTuning()
	tuning.AbandonedOutpostCount = 2
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, nil, 0, tuning, false))

	// Assert
	mainObjectTypes := make([]string, 0, len(zone.MainObjects))
	for _, mainObject := range zone.MainObjects {
		mainObjectTypes = append(mainObjectTypes, mainObject.Type)
	}
	assert.Equal(t, []string{"City", "AbandonedOutpost", "AbandonedOutpost"}, mainObjectTypes)
}

func TestWhenNeutralZoneHasOutpostsAndFootholds_RoadsIncludeBoth(t *testing.T) {
	t.Parallel()
	// Arrange
	plan := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 1}
	tuning := newUnitTuning()
	tuning.AbandonedOutpostCount = 1
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateNeutralZone(newNeutralRequest(plan, []string{"Gate-D"}, 1, tuning, false))

	// Assert
	assert.Equal(t, []template_model.Road{
		{
			Type: "Stone",
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-D"}},
		},
	}, zone.Roads)
}
