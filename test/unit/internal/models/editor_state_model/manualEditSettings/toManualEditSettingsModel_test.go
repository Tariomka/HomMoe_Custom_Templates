package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
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

func TestWhenTheGroupIsPersisted_TheConnectionsAreWrapped(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsManualEditSettings()

	// Act
	model := editor_state_model.ToManualEditSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]editor_state_model.ManualConnectionSave{{ManualConnectionSave: entity.ManualConnections[0]}},
		model.ManualConnections)
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
		ManualZones: []editor_state.ManualZoneSave{{
			Zone:           entities.Zone{Name: "Zone A"},
			ManualPosition: &[2]float64{0.25, 0.75},
		}},
		ManualConnections: []editor_state.ManualConnectionSave{{
			Connection:  entities.Connection{Name: "Conn A"},
			IsUserAdded: true,
		}},
	}
}
