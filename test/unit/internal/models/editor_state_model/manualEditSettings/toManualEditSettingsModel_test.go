package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheGroupIsPersisted_TheZonesAreLifted(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsManualEditSettings()

	// Act
	model := editor_state_model.ToManualEditSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]template_model.Zone{{Name: "Zone A", ManualPosition: new(data.NewVec2(0.25, 0.75))}},
		model.ManualZones)
}

func TestWhenTheGroupIsPersisted_TheConnectionsAreRebuiltFromTheirOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsManualEditSettings()

	// Act
	model := editor_state_model.ToManualEditSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]editor_state_model.ManualConnectionSave{{
			Connection: template_entity.Connection{
				Name: "Conn A", From: "Zone A", To: "Zone B", ConnectionType: "Direct", GuardValue: 4000,
			},
			IsUserAdded: true}},
		model.ManualConnections)
}

func TestWhenAZoneFieldIsUnset_ItIsRebuiltAsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ManualEditSettings{Zones: []editor_state.ManualZoneEdit{{
		Name:           "Zone A",
		FieldOverrides: editor_state.ZoneFieldOverrides{GuardedContentPool: []string{}},
		Unset:          []string{"roads"},
	}}}

	// Act
	model := editor_state_model.ToManualEditSettingsModel(entity)

	// Assert
	assert.Equal(t,
		[]any{[]string{}, true},
		[]any{model.ManualZones[0].GuardedContentPool, model.ManualZones[0].Roads == nil})
}

func TestWhenAnEmptyGroupIsPersisted_TheListsStayNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ManualEditSettings{}

	// Act
	model := editor_state_model.ToManualEditSettingsModel(entity)

	// Assert
	assert.Equal(t, editor_state_model.ManualEditSettings{}, model)
}

func allFieldsManualEditSettings() editor_state.ManualEditSettings {
	return editor_state.ManualEditSettings{
		Zones: []editor_state.ManualZoneEdit{{
			Name:           "Zone A",
			ManualPosition: new(data.NewVec2(0.25, 0.75)),
		}},
		Connections: []editor_state.ManualConnectionEdit{{
			Name:           "Conn A",
			From:           "Zone A",
			To:             "Zone B",
			ConnectionType: "Direct",
			IsUserAdded:    true,
			FieldOverrides: editor_state.ConnectionFieldOverrides{GuardValue: new(4000)},
		}},
	}
}
