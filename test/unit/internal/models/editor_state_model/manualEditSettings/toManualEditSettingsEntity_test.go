package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheGroupIsUnwrapped_TheEntityMatchesTheOriginal(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsManualEditSettings()

	// Act
	unwrapped := editor_state_model.ToManualEditSettingsEntity(
		editor_state_model.ToManualEditSettingsModel(entity))

	// Assert
	assert.Equal(t, entity, unwrapped)
}

func TestWhenAnEmptyGroupIsUnwrapped_TheEntityIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	model := editor_state_model.ManualEditSettings{}

	// Act
	entity := editor_state_model.ToManualEditSettingsEntity(model)

	// Assert
	assert.Equal(t, editor_state.ManualEditSettings{}, entity)
}
