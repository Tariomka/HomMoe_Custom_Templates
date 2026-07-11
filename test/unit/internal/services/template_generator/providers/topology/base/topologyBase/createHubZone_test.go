package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenHubZoneIsCreated_NameIsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Equal(t, "Hub", zone.Name)
}

func TestWhenHoldCityHubHasNoCastles_ForcesSingleCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), true, 1.0, 0, true, "")

	// Assert
	assert.Len(t, zone.MainObjects, 1)
}

func TestWhenHubIsHoldCity_PrimaryCastleCarriesHoldCityWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), true, 1.0, 1, true, "")

	// Assert
	assert.True(t, zone.MainObjects[0].HoldCityWinCon)
}

func TestWhenMandatoryContentNameProvided_ZoneReferencesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), false, 1.0, 1, true, "mandatory_content_hub")

	// Assert
	assert.Contains(t, zone.MandatoryContent, "mandatory_content_hub")
}

func TestWhenMandatoryContentNameIsEmpty_ZoneHasNoMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), false, 1.0, 1, true, "")

	// Assert
	assert.Empty(t, zone.MandatoryContent)
}

func TestWhenHubHasNoCastlesAndIsNotHoldCity_BiomeMatchesZone(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), false, 1.0, 0, true, "")

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "MatchZone"}, zone.ZoneBiome)
}

func TestWhenHubHasMultipleCastles_MainObjectCountMatchesCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateHubZone(nil, newUnitTuning(), false, 1.0, 3, true, "")

	// Assert
	assert.Len(t, zone.MainObjects, 3)
}
