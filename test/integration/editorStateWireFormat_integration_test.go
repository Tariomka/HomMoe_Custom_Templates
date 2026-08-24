package integration_test

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// persistedEditorStateFieldCount is the number of settings keys the editor state
// writes to a .gen.json file, excluding the schema version.
const persistedEditorStateFieldCount = 72

const schemaVersionKey = "schemaVersion"

func TestWhenTheLegacyStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	data := readEditorStateFixture(t, "editorState_v0_flat.gen.json")
	var loaded editor_state.EditorState
	// The legacy file predates the version key, so it decodes at 0; the mapper
	// is what lands it at the current version.
	expected := test_helpers.NewAllFieldsEditorStateEntity()
	expected.SchemaVersion = 0

	// Act
	err := json.Unmarshal(data, &loaded)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestWhenTheLegacyStateFixtureIsMappedToAState_ItLandsAtTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	var loaded editor_state.EditorState
	require.NoError(t, json.Unmarshal(readEditorStateFixture(t, "editorState_v0_flat.gen.json"), &loaded))

	// Act
	state := mappers.NewEditorStateEntityMapper().ToModel(loaded)

	// Assert
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, state.SchemaVersion)
}

func TestWhenTheLegacyStateFixtureIsParsed_ItCarriesEveryPersistedKey(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	keys := parseEditorStateFixtureKeys(t, "editorState_v0_flat.gen.json")

	// Assert
	assert.Len(t, keys, persistedEditorStateFieldCount)
}

func TestWhenTheLegacyStateFixtureIsParsed_ItCarriesNoSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	keys := parseEditorStateFixtureKeys(t, "editorState_v0_flat.gen.json")

	// Assert
	assert.NotContains(t, keys, schemaVersionKey)
}

func TestWhenTheCurrentStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	data := readEditorStateFixture(t, "editorState_v1_flat.gen.json")
	var loaded editor_state.EditorState

	// Act
	err := json.Unmarshal(data, &loaded)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, test_helpers.NewAllFieldsEditorStateEntity(), loaded)
}

func TestWhenBothStateFixturesAreParsed_TheyShareEveryKeyButTheSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	legacyKeys := slices.Sorted(maps.Keys(parseEditorStateFixtureKeys(t, "editorState_v0_flat.gen.json")))
	currentKeys := parseEditorStateFixtureKeys(t, "editorState_v1_flat.gen.json")

	// Act
	delete(currentKeys, schemaVersionKey)

	// Assert
	assert.Equal(t, legacyKeys, slices.Sorted(maps.Keys(currentKeys)))
}

func TestWhenTheAllFieldsStateIsWritten_ItMatchesTheCurrentFixture(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := map[string]any{}
	require.NoError(t, json.Unmarshal(readEditorStateFixture(t, "editorState_v1_flat.gen.json"), &expected))

	// Act
	written, err := json.Marshal(test_helpers.NewAllFieldsEditorStateEntity())

	// Assert
	require.NoError(t, err)
	actual := map[string]any{}
	require.NoError(t, json.Unmarshal(written, &actual))
	assert.Equal(t, expected, actual)
}

// readEditorStateFixture returns a committed .gen.json fixture. The `_v0_` one
// freezes the wire format the editor produced before the state was split into
// entity groups and versioned; the `_v1_` one is what the current writer emits.
// Regrouping the fields reorders the keys, so every assertion above compares
// parsed values and never raw bytes.
func readEditorStateFixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "test_helpers", "testdata", name))
	require.NoError(t, err)

	return data
}

func parseEditorStateFixtureKeys(t *testing.T, name string) map[string]json.RawMessage {
	t.Helper()

	keys := map[string]json.RawMessage{}
	require.NoError(t, json.Unmarshal(readEditorStateFixture(t, name), &keys))

	return keys
}
