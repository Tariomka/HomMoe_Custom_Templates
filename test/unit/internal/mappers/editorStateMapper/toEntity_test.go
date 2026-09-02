package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

// The check is a round trip rather than a field-by-field expectation so that a
// group the mapper forgets to carry fails the test instead of being mirrored by
// an expectation written from the same list. It cannot go through JSON any
// more: the model re-declares the content and manual-edit groups to hold model
// element types, so it no longer carries the entity's json tags.
func TestWhenAStateIsMappedToTheEntityAndBack_EveryPersistedFieldSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()
	mapper := mappers.NewEditorStateMapper()

	// Act
	roundTripped := mapper.ToModel(mapper.ToEntity(state))

	// Assert
	assert.Equal(t, state, roundTripped)
}

func TestWhenAStateIsMappedToTheEntity_ItCarriesTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()

	// Act
	entity := mappers.NewEditorStateMapper().ToEntity(state)

	// Assert
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, entity.SchemaVersion)
}
