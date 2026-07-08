package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsCreated_NameCombinesNeutralPrefixWithPlanLabel(t *testing.T) {
	// Arrange
	label := gofakeit.LetterN(3)
	plan := models.NeutralZonePlan{Label: label, Quality: models.QualityMedium, CastleCount: 1}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), false)

	// Assert
	assert.Equal(t, "Neutral-"+label, zone.Name)
}

func TestWhenPlanHasCastles_MainObjectCountMatchesCastleCount(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityMedium, CastleCount: 2}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), false)

	// Assert
	assert.Len(t, zone.MainObjects, 2)
}

func TestWhenHoldCityZoneHasNoPlannedCastles_ForcesSingleCastle(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityMedium, CastleCount: 0}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), true)

	// Assert
	assert.Len(t, zone.MainObjects, 1)
}

func TestWhenZoneIsHoldCity_PrimaryCastleCarriesHoldCityWinCondition(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityHigh, CastleCount: 1}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), true)

	// Assert
	assert.True(t, zone.MainObjects[0].HoldCityWinCon)
}

func TestWhenZoneHasNoMainObjects_BiomeMatchesZone(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityLow, CastleCount: 0}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), false)

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "MatchZone"}, zone.ZoneBiome)
}

func TestWhenZoneHasMainObjects_BiomeMatchesPrimaryMainObject(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityLow, CastleCount: 1}
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, newUnitTuning(), false)

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}, zone.ZoneBiome)
}

func TestWhenAbandonedOutpostCountIsPositive_AppendsOutpostsAfterCastles(t *testing.T) {
	// Arrange
	plan := models.NeutralZonePlan{Label: "D", Quality: models.QualityMedium, CastleCount: 1}
	tuning := newUnitTuning()
	tuning.AbandonedOutpostCount = 2
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateNeutralZone(plan, nil, 1.0, 0, true, tuning, false)

	// Assert
	mainObjectTypes := make([]string, 0, len(zone.MainObjects))
	for _, mainObject := range zone.MainObjects {
		mainObjectTypes = append(mainObjectTypes, mainObject.Type)
	}
	assert.Equal(t, []string{"City", "AbandonedOutpost", "AbandonedOutpost"}, mainObjectTypes)
}
