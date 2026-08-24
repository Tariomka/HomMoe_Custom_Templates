package editorStateEntity_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheSchemaVersionIsUnset_TheWrittenStateCarriesTheCurrentVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, "1", string(keys["schemaVersion"]))
}

func TestWhenTheSchemaVersionIsUnexpected_TheWrittenStateStillCarriesTheCurrentVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{SchemaVersion: 99}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, "1", string(keys["schemaVersion"]))
}

func TestWhenTheStateIsWritten_TheGroupFieldsSitBesideTheSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{TemplateName: "Flat"}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, `"Flat"`, string(keys["templateName"]))
}

// The two tests below characterise the `omitempty` tags rather than endorse
// them: a nil slice and an empty one are indistinguishable once written, so the
// distinction EqualsIgnoringManualEdits draws between them exists in memory
// only. Loading either back yields whatever the reader seeded.
func TestWhenContentRowsAreNil_TheKeyIsNotWritten(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{PlayerZoneContentRows: nil}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, keys, "playerZoneContentRows")
}

func TestWhenContentRowsAreEmpty_TheKeyIsNotWrittenEither(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{PlayerZoneContentRows: []editor_state.ZoneContentRow{}}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, keys, "playerZoneContentRows")
}

func marshalToKeys(t *testing.T, entity editor_state.EditorStateEntity) map[string]json.RawMessage {
	t.Helper()

	data, err := json.Marshal(entity)
	require.NoError(t, err)

	keys := map[string]json.RawMessage{}
	require.NoError(t, json.Unmarshal(data, &keys))

	return keys
}
