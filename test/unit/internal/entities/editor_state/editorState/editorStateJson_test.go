package editorState_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheStateIsWritten_TheGroupFieldsSitBesideTheSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{TemplateName: "Flat", SchemaVersion: 1}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, `"Flat"`, string(keys["templateName"]))
}

// The two tests below characterize the `omitempty` tags rather than endorse
// them: a nil slice and an empty one are indistinguishable once written, so the
// distinction EqualsIgnoringManualEdits draws between them exists in memory
// only. Loading either back yields whatever the reader seeded.
func TestWhenContentRowsAreNil_TheKeyIsNotWritten(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{PlayerZoneContentRows: nil}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, keys, "playerZoneContentRows")
}

func TestWhenContentRowsAreEmpty_TheKeyIsNotWrittenEither(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{PlayerZoneContentRows: []editor_state.ZoneContentRow{}}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, keys, "playerZoneContentRows")
}

// The load path seeds the entity with the defaults before decoding, so this
// merge is what lets a key the file omits keep its default instead of
// collapsing to a zero value.
func TestWhenAKeyIsAbsentFromTheFile_TheSeededValueSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := json.Unmarshal([]byte(`{"playerCount":5}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", entity.TemplateName)
}

func TestWhenAKeyIsPresentInTheFile_ItReplacesTheSeededValue(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := json.Unmarshal([]byte(`{"templateName":"From File"}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "From File", entity.TemplateName)
}

func TestWhenTheFileCarriesNoSchemaVersion_TheStateDecodesAtVersionZero(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{}

	// Act
	err := json.Unmarshal([]byte(`{"templateName":"Legacy"}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 0, entity.SchemaVersion)
}

func marshalToKeys(t *testing.T, entity editor_state.EditorState) map[string]jsontext.Value {
	t.Helper()

	data, err := json.Marshal(entity)
	require.NoError(t, err)

	keys := map[string]jsontext.Value{}
	require.NoError(t, json.Unmarshal(data, &keys))

	return keys
}
