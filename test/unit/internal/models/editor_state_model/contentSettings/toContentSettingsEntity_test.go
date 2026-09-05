package contentSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

// The round trip is the drift guard: a member the converter drops in either
// direction cannot survive it, and unlike a hand-written expectation it is not
// derived from the same field list the converter uses.
func TestWhenTheGroupIsUnwrappedAndWrappedAgain_EveryFieldSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	model := editor_state_model.ToContentSettingsModel(allFieldsContentSettings())

	// Act
	roundTripped := editor_state_model.ToContentSettingsModel(
		editor_state_model.ToContentSettingsEntity(model))

	// Assert
	assert.Equal(t, model, roundTripped)
}

func TestWhenTheGroupIsUnwrapped_TheEntityMatchesTheOriginal(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsContentSettings()

	// Act
	unwrapped := editor_state_model.ToContentSettingsEntity(
		editor_state_model.ToContentSettingsModel(entity))

	// Assert
	assert.Equal(t, entity, unwrapped)
}

func TestWhenAnEmptyGroupIsUnwrapped_TheEntityIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	model := editor_state_model.ContentSettings{}

	// Act
	entity := editor_state_model.ToContentSettingsEntity(model)

	// Assert
	assert.Equal(t, editor_state.ContentSettings{}, entity)
}
