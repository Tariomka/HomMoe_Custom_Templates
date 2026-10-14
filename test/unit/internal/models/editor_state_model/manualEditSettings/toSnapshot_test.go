package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenRichZonesAreRoundTripped_TheZonesAreReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := newRichZones()
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, newRichConnections())

	// Act
	decodedZones, _ := settings.ToSnapshot()

	// Assert
	assert.Equal(t, zones, decodedZones)
}

func TestWhenRichConnectionsAreRoundTripped_TheConnectionsAreReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := newRichConnections()
	settings := editor_state_model.ManualEditSettingsFromSnapshot(newRichZones(), connections)

	// Act
	_, decodedConnections := settings.ToSnapshot()

	// Assert
	assert.Equal(t, connections, decodedConnections)
}

func TestWhenTheSettingsAreEmpty_DecodedZonesAreNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var settings editor_state_model.ManualEditSettings

	// Act
	zones, _ := settings.ToSnapshot()

	// Assert
	assert.Nil(t, zones)
}

func TestWhenTheSettingsAreEmpty_DecodedConnectionsAreNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var settings editor_state_model.ManualEditSettings

	// Act
	_, connections := settings.ToSnapshot()

	// Assert
	assert.Nil(t, connections)
}

func TestWhenTheDecodedZoneIsMutated_ReDecodingYieldsTheOriginalZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := newRichZones()
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)
	decoded, _ := settings.ToSnapshot()
	decoded[0].GuardedContentPool[0] = "mutated"

	// Act
	redecoded, _ := settings.ToSnapshot()

	// Assert
	assert.Equal(t, zones[0].GuardedContentPool, redecoded[0].GuardedContentPool)
}

func TestWhenTheDecodedConnectionIsMutated_ReDecodingYieldsTheOriginalConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := newRichConnections()
	settings := editor_state_model.ManualEditSettingsFromSnapshot(nil, connections)
	_, decoded := settings.ToSnapshot()
	decoded[0].PortalPlacementRulesFrom[0].Type = "mutated"

	// Act
	_, redecoded := settings.ToSnapshot()

	// Assert
	assert.Equal(t, "rule", redecoded[0].PortalPlacementRulesFrom[0].Type)
}

func TestWhenTheSettingsHoldEditsAgainstADerivation_SnapshotDecodesThemAgainstAnEmptyZone(t *testing.T) {
	t.Parallel()
	// Arrange
	var settings editor_state_model.ManualEditSettings
	settings.ManualZones = []editor_state_model.ManualZoneEdit{{}}
	settings.ManualZones[0].Name = "Only"
	settings.ManualZones[0].FieldOverrides.Size = new(3.0)

	// Act
	zones, _ := settings.ToSnapshot()

	// Assert
	assert.InDelta(t, 3.0, zones[0].Size, 0)
}
