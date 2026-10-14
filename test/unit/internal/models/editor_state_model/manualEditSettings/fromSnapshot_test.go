package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSnapshotHasZones_OneEditIsCreatedPerZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := newRichZones()

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)

	// Assert
	assert.Len(t, settings.ManualZones, len(zones))
}

func TestWhenTheSnapshotHasConnections_OneEditIsCreatedPerConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := newRichConnections()

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(nil, connections)

	// Assert
	assert.Len(t, settings.ManualConnections, len(connections))
}

func TestWhenTheSnapshotHasZones_EditsKeepTheZoneOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := newRichZones()

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)

	// Assert
	assert.Equal(t, []string{zones[0].Name, zones[1].Name},
		[]string{settings.ManualZones[0].Name, settings.ManualZones[1].Name})
}

func TestWhenTheSnapshotIsNil_ZoneListStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(nil, nil)

	// Assert
	assert.Nil(t, settings.ManualZones)
}

func TestWhenTheSnapshotIsNil_ConnectionListStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(nil, nil)

	// Assert
	assert.Nil(t, settings.ManualConnections)
}

func TestWhenTheSnapshotIsEmpty_ZoneListIsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{}

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)

	// Assert
	assert.Nil(t, settings.ManualZones)
}

func TestWhenAZoneHoldsNoData_ItsEditCarriesNoOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Empty"}}

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)

	// Assert
	assert.Equal(t, editor_state.ZoneFieldOverrides{}, settings.ManualZones[0].FieldOverrides)
}

func TestWhenAZoneHoldsAField_ItsEditStoresTheFieldAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Sized", Size: 2.5}}

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(zones, nil)

	// Assert
	assert.Equal(t, new(2.5), settings.ManualZones[0].FieldOverrides.Size)
}

func TestWhenAConnectionHoldsAField_ItsEditStoresTheFieldAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{{Name: "A-B", Length: 4.5}}

	// Act
	settings := editor_state_model.ManualEditSettingsFromSnapshot(nil, connections)

	// Assert
	assert.Equal(t, new(4.5), settings.ManualConnections[0].FieldOverrides.Length)
}
