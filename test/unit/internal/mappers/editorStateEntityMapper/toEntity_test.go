package editorStateEntityMapper_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The comparison runs through JSON rather than through the nine group fields so
// that a group the mapper forgets to carry fails the test instead of being
// mirrored by an expectation written from the same list.
func TestWhenAStateIsMappedToTheEntity_EveryPersistedFieldIsCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()
	expected := map[string]any{}
	stateJSON, err := json.Marshal(state)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(stateJSON, &expected))

	// Act
	entity := mappers.NewEditorStateEntityMapper().ToEntity(state)

	// Assert
	actual := map[string]any{}
	entityJSON, marshalErr := json.Marshal(entity)
	require.NoError(t, marshalErr)
	require.NoError(t, json.Unmarshal(entityJSON, &actual))
	delete(actual, "schemaVersion")
	assert.Equal(t, expected, actual)
}

func TestWhenAStateIsMappedToTheEntity_ItCarriesTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()

	// Act
	entity := mappers.NewEditorStateEntityMapper().ToEntity(state)

	// Assert
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, entity.SchemaVersion)
}
