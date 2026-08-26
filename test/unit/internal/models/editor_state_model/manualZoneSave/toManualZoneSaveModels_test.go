package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoPersistedZoneSaves_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state.ManualZoneSave

	// Act
	models := editor_state_model.ToManualZoneSaveModels(saves)

	// Assert
	assert.Nil(t, models)
}

func TestWhenZoneSavesArePersisted_EachOneIsWrapped(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualZoneSave{{
		Zone:           entities.Zone{Name: "Zone A"},
		ManualPosition: &[2]float64{0.25, 0.75},
	}}

	// Act
	models := editor_state_model.ToManualZoneSaveModels(saves)

	// Assert
	assert.Equal(t, []editor_state_model.ManualZoneSave{{ManualZoneSave: saves[0]}}, models)
}

func TestWhenThereAreNoZoneSaveModels_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state_model.ManualZoneSave

	// Act
	unwrapped := editor_state_model.ToManualZoneSaveEntities(saves)

	// Assert
	assert.Nil(t, unwrapped)
}

func TestWhenZoneSavesAreUnwrapped_TheEntitiesAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state_model.ManualZoneSave{{
		Zone:           entities.Zone{Name: "Zone A"},
		ManualPosition: &[2]float64{0.25, 0.75},
	}}

	// Act
	unwrapped := editor_state_model.ToManualZoneSaveEntities(saves)

	// Assert
	assert.Equal(t, []editor_state.ManualZoneSave{saves[0].ManualZoneSave}, unwrapped)
}
