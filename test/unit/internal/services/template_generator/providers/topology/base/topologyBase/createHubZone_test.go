package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenHubZoneIsCreated_NameIsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(t, "Hub", zone.Name)
}

func TestWhenHoldCityHubHasNoCastles_ForcesSingleCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), true, 1.0, 0, true, "")

	// Assert
	assert.Len(t, zone.MainObjects, 1)
}

func TestWhenHubIsHoldCity_PrimaryCastleCarriesHoldCityWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), true, 1.0, 1, true, "")

	// Assert
	assert.True(t, zone.MainObjects[0].HoldCityWinCon)
}

func TestWhenMandatoryContentNameProvided_ZoneReferencesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone(
		"Hub", nil, newUnitTuning(), false, 1.0, 1, true, "mandatory_content_hub")

	// Assert
	assert.Contains(t, zone.MandatoryContent, "mandatory_content_hub")
}

func TestWhenMandatoryContentNameIsEmpty_ZoneHasNoMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Empty(t, zone.MandatoryContent)
}

func TestWhenHubHasNoCastlesAndIsNotHoldCity_BiomeMatchesZone(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 0, true, "")

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "MatchZone"}, zone.ZoneBiome)
}

func TestWhenHubHasMultipleCastles_MainObjectCountMatchesCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 3, true, "")

	// Assert
	assert.Len(t, zone.MainObjects, 3)
}

func TestWhenHubZoneIsCreated_UsesHighestProfileGuardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(
		t,
		common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest).GuardedContentPool,
		zone.GuardedContentPool,
	)
}

func TestWhenHubZoneIsCreated_UsesHighestProfileResourcesPool(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(
		t,
		common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest).ResourcesContentPool,
		zone.ResourcesContentPool,
	)
}

func TestWhenHubZoneIsCreated_UsesHighestProfileGuardedContentValue(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(t, 960000, zone.GuardedContentValue)
}

func TestWhenHubZoneIsCreated_UsesHighestProfileGuardMultiplier(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.InDelta(t, 2.3, zone.GuardMultiplier, 1e-9)
}

func TestWhenHubZoneIsCreated_ClassifiesAsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, zone_services.NewZoneClassifier().GetQuality(zone))
}

func TestWhenHubZoneIsCreated_RoadsCountCastlesOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	tuning := newUnitTuning()
	tuning.AbandonedOutpostCount = 2
	tuning.RemoteFootholdCount = 2
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateHubZone("Hub", []string{"Hub-Gate"}, tuning, false, 1.0, 2, true, "")

	// Assert
	assert.Equal(t, []entities.Road{
		{
			Type: "Stone",
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "Connection", Args: []string{"Hub-Gate"}},
		},
	}, zone.Roads)
}
