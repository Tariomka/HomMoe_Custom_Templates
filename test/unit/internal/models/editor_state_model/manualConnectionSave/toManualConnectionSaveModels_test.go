package manualConnectionSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoPersistedConnectionSaves_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state.ManualConnectionSave

	// Act
	models := editor_state_model.ToManualConnectionSaveModels(saves)

	// Assert
	assert.Nil(t, models)
}

func TestWhenConnectionSavesArePersisted_EachOneIsWrapped(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualConnectionSave{{
		Connection:  entities.Connection{Name: "Conn A"},
		IsUserAdded: true,
	}}

	// Act
	models := editor_state_model.ToManualConnectionSaveModels(saves)

	// Assert
	assert.Equal(t, []editor_state_model.ManualConnectionSave{{ManualConnectionSave: saves[0]}}, models)
}

func TestWhenThereAreNoConnectionSaveModels_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state_model.ManualConnectionSave

	// Act
	unwrapped := editor_state_model.ToManualConnectionSaveEntities(saves)

	// Assert
	assert.Nil(t, unwrapped)
}

func TestWhenConnectionSavesAreUnwrapped_TheEntitiesAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state_model.ManualConnectionSave{{
		Connection:  entities.Connection{Name: "Conn A"},
		IsUserAdded: true,
	}}

	// Act
	unwrapped := editor_state_model.ToManualConnectionSaveEntities(saves)

	// Assert
	assert.Equal(t, []editor_state.ManualConnectionSave{saves[0].ManualConnectionSave}, unwrapped)
}
