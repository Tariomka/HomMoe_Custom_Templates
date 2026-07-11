package zoneLayoutProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreatingZoneLayouts_ReturnsFourLayouts(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewZoneLayoutProvider()

	// Act
	layouts := provider.CreateZoneLayouts()

	// Assert
	assert.Len(t, layouts, 4)
}

func TestWhenCreatingZoneLayouts_NamesLayoutsAfterRegistryLayoutSids(t *testing.T) {
	t.Parallel()
	// Arrange
	layoutSids := registry.GetLayoutValues()
	provider := providers.NewZoneLayoutProvider()

	// Act
	layouts := provider.CreateZoneLayouts()

	// Assert
	var actualNames []string
	for _, layout := range layouts {
		actualNames = append(actualNames, layout.Name)
	}
	assert.Equal(t, []string{
		layoutSids.Spawns,
		layoutSids.Sides,
		layoutSids.TreasureZone,
		layoutSids.Center,
	}, actualNames)
}

func TestWhenCreatingZoneLayouts_SpawnsLayoutMatchesExpectedTuning(t *testing.T) {
	t.Parallel()
	// Arrange
	layoutSids := registry.GetLayoutValues()
	expected := entities.ZoneLayoutDef{
		Name:                  layoutSids.Spawns,
		ObstaclesFill:         0.24,
		ObstaclesFillVoid:     0.48,
		LakesFill:             0.30,
		MinLakeArea:           16,
		ElevationClusterScale: 0.16,
		ElevationModes: []entities.ElevationMode{
			{Weight: 2, MinElevatedFraction: 0.2, MaxElevatedFraction: 0.4},
			{Weight: 1, MinElevatedFraction: 0.6, MaxElevatedFraction: 0.8},
		},
		RoadClusterArea: 160,
		GuardedEncounterResourceFractions: entities.GuardedEncounterResourceFractions{
			CountBounds: []int{},
			Fractions:   []float64{0.66},
		},
		AmbientPickupDistribution: entities.AmbientPickupDistribution{
			Repulsion:          1.0,
			Noise:              0.4,
			RoadAttraction:     -0.30,
			ObstacleAttraction: 0,
			GroupSizeWeights:   []int{20, 2, 1},
		},
	}
	provider := providers.NewZoneLayoutProvider()

	// Act
	layouts := provider.CreateZoneLayouts()

	// Assert
	assert.Equal(t, expected, layouts[0])
}

func TestWhenCreatingZoneLayouts_EveryLayoutHasTwoElevationModes(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewZoneLayoutProvider()

	// Act
	layouts := provider.CreateZoneLayouts()

	// Assert
	for _, layout := range layouts {
		assert.Len(t, layout.ElevationModes, 2, "layout %q should have two elevation modes", layout.Name)
	}
}
