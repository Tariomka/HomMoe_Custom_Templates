package integration_test

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// persistedEditorStateFieldCount is the number of settings keys the editor state
// writes to a .gen.json file, excluding the schema version.
const persistedEditorStateFieldCount = 72

const schemaVersionKey = "schemaVersion"

// The v0 and v1 files wrote manualPosition as [x, y], which encoding/json/v2
// refuses to read into the current struct, so the migration is the only way in.
// Everything they do carry has to come back untouched.
func TestWhenTheLegacyStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newMigratedAllFieldsEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, "editorState_v0_flat.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestWhenTheLegacyStateFixtureIsLoaded_ItLandsAtTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	loaded, err := decodeEditorStateFixture(t, "editorState_v0_flat.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, loaded.SchemaVersion)
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

func TestWhenTheV1StateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newMigratedAllFieldsEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, "editorState_v1_flat.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestWhenTheCurrentStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := test_helpers.NewAllFieldsEditorStateEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, "editorState_v2_flat.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

// The generator stamps are the one thing v2 adds that no older file can carry,
// so they are also the one thing that proves the new keys really round-trip.
func TestWhenTheCurrentStateFixtureIsLoaded_TheGeneratorStampsComeBack(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := test_helpers.NewAllFieldsEditorStateEntity().ManualZones[0]

	// Act
	loaded, err := decodeEditorStateFixture(t, "editorState_v2_flat.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded.ManualZones[0])
}

func TestWhenAStateFileClaimsANewerSchema_TheLoadIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	bumped := strings.Replace(
		string(readEditorStateFixture(t, "editorState_v2_flat.gen.json")),
		`"schemaVersion": 2`,
		`"schemaVersion": 3`,
		1)
	statePath := filepath.Join(t.TempDir(), "tooNew.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(bumped), 0o644))
	target := mappers.NewEditorStateMapper().NewDefaultEntity()

	// Act
	err := newEditorStateMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedSchemaVersionError))
}

// Files the user already has on disk predate the version key entirely. If any
// of them stops loading, the migration is wrong, not the file. output/ is
// gitignored, so this is an opt-in local check rather than a CI gate.
func TestWhenTheStatesInTheOutputFolderAreLoaded_TheyAllStillLoad(t *testing.T) {
	t.Parallel()
	// Arrange
	paths, err := filepath.Glob(filepath.Join("..", "..", "output", "*.gen.json"))
	require.NoError(t, err)
	if len(paths) == 0 {
		t.Skip("no local .gen.json files in output/ to migrate")
	}
	migrator := newEditorStateMigrator()

	// Act & Assert
	for _, path := range paths {
		target := mappers.NewEditorStateMapper().NewDefaultEntity()
		assert.NoError(t, migrator.Load(path, &target), "%s no longer loads", filepath.Base(path))
	}
}

func TestWhenBothStateFixturesAreParsed_TheyShareEveryKeyButTheSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	legacyKeys := slices.Sorted(maps.Keys(parseEditorStateFixtureKeys(t, "editorState_v0_flat.gen.json")))
	currentKeys := parseEditorStateFixtureKeys(t, "editorState_v2_flat.gen.json")

	// Act
	delete(currentKeys, schemaVersionKey)

	// Assert
	assert.Equal(t, legacyKeys, slices.Sorted(maps.Keys(currentKeys)))
}

func TestWhenTheAllFieldsStateIsWritten_ItMatchesTheCurrentFixture(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := map[string]any{}
	require.NoError(t, json.Unmarshal(readEditorStateFixture(t, "editorState_v2_flat.gen.json"), &expected))

	// Act
	//nolint:musttag // data.Vec2 is deliberately tag-free: its field names X and Y are the wire form.
	written, err := json.Marshal(test_helpers.NewAllFieldsEditorStateEntity())

	// Assert
	require.NoError(t, err)
	actual := map[string]any{}
	require.NoError(t, json.Unmarshal(written, &actual))
	assert.Equal(t, expected, actual)
}

// newMigratedAllFieldsEntity is the all-fields state as it comes back from a
// pre-v2 file: identical but for the generator stamps, which those files had
// nowhere to put.
func newMigratedAllFieldsEntity() editor_state.EditorState {
	expected := test_helpers.NewAllFieldsEditorStateEntity()
	expected.ManualZones[0].GeneratorPosition = nil
	expected.ManualZones[0].GeneratorRing = nil

	return expected
}

// decodeEditorStateFixture reads a fixture the way the app does: seeded with the
// defaults and routed through the migrator, which is the only path that can
// read a pre-v2 file at all.
func decodeEditorStateFixture(t *testing.T, name string) (editor_state.EditorState, error) {
	t.Helper()

	target := mappers.NewEditorStateMapper().NewDefaultEntity()
	err := newEditorStateMigrator().
		Load(filepath.Join("..", "test_helpers", "testdata", name), &target)

	return target, err
}

func newEditorStateMigrator() editor_state_migrator.IEditorStateMigrator {
	return editor_state_migrator.NewEditorStateMigrator(
		repositories.NewEditorStateRepository(),
		repositories.NewLegacyEditorStateRepository())
}

// readEditorStateFixture returns a committed .gen.json fixture. The `_v0_` one
// freezes the wire format the editor produced before the state was split into
// entity groups and versioned; `_v1_` is the last version that wrote
// manualPosition as an array; the `_v2_` one is what the current writer emits.
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
