package gladiatorArenaProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenArenaWinConditionIsDisabled_LeavesZonesUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{{Name: "Hub"}}}
	configuration := *config.NewGeneratorConfig()

	// Act
	newProvider().PlaceArena(configuration, &variant)

	// Assert
	assert.Equal(t, 0, countArenaMainObjects(variant.Zones[0]))
}

func TestWhenVariantHasHubZone_PlacesArenaMainObjectInHub(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{
		newNeutralZone("A", neutral_zone.QualityHighest),
		{Name: "Hub"},
	}}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 1, countArenaMainObjects(variant.Zones[1]))
}

func TestWhenVariantHasHubZone_LeavesNeutralZonesWithoutArena(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{
		newNeutralZone("A", neutral_zone.QualityHighest),
		{Name: "Hub"},
	}}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 0, countArenaMainObjects(variant.Zones[0]))
}

func TestWhenArenaMainObjectIsPlaced_UsesTheUniformPlacementTheGameExpects(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{{Name: "Hub"}}}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t,
		template_model.MainObject{
			Type:          arenaObjectType,
			Placement:     "Uniform",
			PlacementArgs: []string{"true", "0", "0"},
		},
		variant.Zones[0].MainObjects[0])
}

func TestWhenNeutralZonesAreConnected_MarksTheConnectionAsArena(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			newNeutralZone("A", neutral_zone.QualityHigh),
			newNeutralZone("B", neutral_zone.QualityHigh),
		},
		Connections: []template_model.Connection{
			{Name: "Neutral-A-Neutral-B", From: "Neutral-A", To: "Neutral-B"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, arenaObjectType, variant.Connections[0].ConnectionType)
}

func TestWhenNeutralZonesAreConnected_PlacesNoArenaMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			newNeutralZone("A", neutral_zone.QualityHigh),
			newNeutralZone("B", neutral_zone.QualityHigh),
		},
		Connections: []template_model.Connection{
			{Name: "Neutral-A-Neutral-B", From: "Neutral-A", To: "Neutral-B"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 0, countArenaMainObjects(variant.Zones[0])+countArenaMainObjects(variant.Zones[1]))
}

func TestWhenSeveralNeutralConnectionsExist_PicksTheRichestPair(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			newNeutralZone("A", neutral_zone.QualityLowest),
			newNeutralZone("B", neutral_zone.QualityLowest),
			newNeutralZone("C", neutral_zone.QualityHighest),
			newNeutralZone("D", neutral_zone.QualityHighest),
		},
		Connections: []template_model.Connection{
			{Name: "Poor", From: "Neutral-A", To: "Neutral-B"},
			{Name: "Rich", From: "Neutral-C", To: "Neutral-D"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, arenaObjectType, variant.Connections[1].ConnectionType)
}

func TestWhenEquallyRichNeutralConnectionsExist_PicksTheLowestConnectionName(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			newNeutralZone("A", neutral_zone.QualityHigh),
			newNeutralZone("B", neutral_zone.QualityHigh),
			newNeutralZone("C", neutral_zone.QualityHigh),
		},
		Connections: []template_model.Connection{
			{Name: "Zulu", From: "Neutral-B", To: "Neutral-C"},
			{Name: "Alpha", From: "Neutral-A", To: "Neutral-B"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, arenaObjectType, variant.Connections[1].ConnectionType)
}

func TestWhenOnlyPlayerToNeutralConnectionsExist_FallsBackToTheRichestNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			{Name: "Spawn-A"},
			newNeutralZone("B", neutral_zone.QualityLow),
			newNeutralZone("C", neutral_zone.QualityHighest),
		},
		Connections: []template_model.Connection{
			{Name: "Spawn-A-Neutral-B", From: "Spawn-A", To: "Neutral-B"},
			{Name: "Spawn-A-Neutral-C", From: "Spawn-A", To: "Neutral-C"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 1, countArenaMainObjects(variant.Zones[2]))
}

func TestWhenOnlyPlayerToNeutralConnectionsExist_LeavesEveryConnectionUnmarked(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			{Name: "Spawn-A"},
			newNeutralZone("B", neutral_zone.QualityLow),
		},
		Connections: []template_model.Connection{
			{Name: "Spawn-A-Neutral-B", From: "Spawn-A", To: "Neutral-B"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 0, countArenaConnections(variant))
}

func TestWhenEquallyRichNeutralZonesExist_PicksTheLowestZoneName(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{
		newNeutralZone("Z", neutral_zone.QualityHigh),
		newNeutralZone("A", neutral_zone.QualityHigh),
	}}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 1, countArenaMainObjects(variant.Zones[1]))
}

func TestWhenVariantHasNoNeutralZoneAtAll_PlacesNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones:       []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}},
		Connections: []template_model.Connection{{Name: "Spawn-A-Spawn-B", From: "Spawn-A", To: "Spawn-B"}},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t,
		0,
		countArenaConnections(variant)+countArenaMainObjects(variant.Zones[0])+countArenaMainObjects(variant.Zones[1]))
}

// A zone whose content pools cannot be classified scores Unknown, which is -1,
// so inference alone could never let it win the arena. The tier the generator
// planned for it can.
func TestWhenTheRichestZoneCannotBeInferred_TheRecordedTierWinsTheArena(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{Zones: []template_model.Zone{
		newNeutralZone("A", neutral_zone.QualityLow),
		{Name: "Neutral-B", Quality: new(neutral_zone.QualityHighest)},
	}}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, 1, countArenaMainObjects(variant.Zones[1]))
}

func TestWhenConnectionEndpointsCannotBeInferred_TheRecordedTiersPickTheConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	variant := template_model.Variant{
		Zones: []template_model.Zone{
			newNeutralZone("A", neutral_zone.QualityLow),
			newNeutralZone("B", neutral_zone.QualityLow),
			{Name: "Neutral-C", Quality: new(neutral_zone.QualityHighest)},
			{Name: "Neutral-D", Quality: new(neutral_zone.QualityHighest)},
		},
		Connections: []template_model.Connection{
			{Name: "Alpha", From: "Neutral-A", To: "Neutral-B"},
			{Name: "Zulu", From: "Neutral-C", To: "Neutral-D"},
		},
	}

	// Act
	newProvider().PlaceArena(newArenaConfiguration(), &variant)

	// Assert
	assert.Equal(t, arenaObjectType, variant.Connections[1].ConnectionType)
}
