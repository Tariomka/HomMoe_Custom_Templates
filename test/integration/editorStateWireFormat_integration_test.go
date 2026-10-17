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

// legacyEditorStateFieldCount is the number of settings keys the flat pre-v3
// editor state wrote to a .gen.json file, excluding the schema version.
const legacyEditorStateFieldCount = 72

const (
	schemaVersionKey      = "schemaVersion"
	currentStateFixture   = "editorState_v3_flat.gen.json"
	previousStateFixture  = "editorState_v2_flat.gen.json"
	unversionedFixture    = "editorState_v0_flat.gen.json"
	arrayPositionsFixture = "editorState_v1_flat.gen.json"
)

// The v0 and v1 files wrote manualPosition as [x, y], which encoding/json/v2
// refuses to read into a newer struct, so the migration is the only way in.
// Everything they do carry has to come back untouched.
func TestWhenTheUnversionedStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newStampFreeAllFieldsEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, unversionedFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestWhenTheUnversionedStateFixtureIsLoaded_ItLandsAtTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	loaded, err := decodeEditorStateFixture(t, unversionedFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, loaded.SchemaVersion)
}

func TestWhenTheUnversionedStateFixtureIsParsed_ItCarriesEveryFlatKey(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	keys := parseEditorStateFixtureKeys(t, unversionedFixture)

	// Assert
	assert.Len(t, keys, legacyEditorStateFieldCount)
}

func TestWhenTheUnversionedStateFixtureIsParsed_ItCarriesNoSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	keys := parseEditorStateFixtureKeys(t, unversionedFixture)

	// Assert
	assert.NotContains(t, keys, schemaVersionKey)
}

func TestWhenTheV1StateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newStampFreeAllFieldsEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, arrayPositionsFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

// v2 is the last flat format: its ban and override text, its full snapshot and
// its flat groups all have to land in the sectioned v3 state unchanged.
func TestWhenTheV2StateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := test_helpers.NewAllFieldsEditorStateEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, previousStateFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

func TestWhenTheCurrentStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := test_helpers.NewAllFieldsEditorStateEntity()

	// Act
	loaded, err := decodeEditorStateFixture(t, currentStateFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded)
}

// The generator stamps are what no pre-v2 file can carry, so they prove the
// structure keys really round-trip.
func TestWhenTheCurrentStateFixtureIsLoaded_TheGeneratorStampsComeBack(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := test_helpers.NewAllFieldsEditorStateEntity().ManualEdits.Zones[0]

	// Act
	loaded, err := decodeEditorStateFixture(t, currentStateFixture)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, loaded.ManualEdits.Zones[0])
}

func TestWhenTheCurrentStateFixtureIsParsed_OnlyIdentityMapAndSchemaFieldsSitAtTheRoot(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	keys := parseEditorStateFixtureKeys(t, currentStateFixture)

	// Assert
	assert.Equal(t,
		[]string{
			"castles", "content", "experimentalMapSizes", "gameMode", "gameRules", "generation", "generatorVersion",
			"manualEdits", "mapSize", "neutralZones", "players", "schemaVersion", "templateName",
		},
		slices.Sorted(maps.Keys(keys)))
}

// Sectioning moves keys, it never renames them: every flat v2 settings key is
// found again inside exactly one v3 section or at the root.
func TestWhenTheCurrentStateFixtureIsParsed_EveryFlatSettingsKeyReappearsOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	flatKeys := parseEditorStateFixtureKeys(t, previousStateFixture)
	delete(flatKeys, "manualZones")
	delete(flatKeys, "manualConnections")
	currentKeys := parseEditorStateFixtureKeys(t, currentStateFixture)

	// Act
	var found []string
	for key, value := range currentKeys {
		if !slices.Contains([]string{"players", "neutralZones", "castles", "generation", "gameRules", "content"}, key) {
			found = append(found, key)
			continue
		}

		section := map[string]json.RawMessage{}
		require.NoError(t, json.Unmarshal(value, &section))
		found = slices.AppendSeq(found, maps.Keys(section))
	}

	// Assert
	assert.Subset(t, found, slices.Collect(maps.Keys(flatKeys)))
}

func TestWhenAStateFileClaimsANewerSchema_TheLoadIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	bumped := strings.Replace(
		string(readEditorStateFixture(t, currentStateFixture)),
		`"schemaVersion": 3`,
		`"schemaVersion": 4`,
		1)
	statePath := filepath.Join(t.TempDir(), "tooNew.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(bumped), 0o644))
	target := mappers.NewEditorStateMapper().NewDefaultEntity()

	// Act
	_, err := newEditorStateMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedSchemaVersionError))
}

func TestWhenAStateFileClaimsANewerContentDefaultsVersion_TheLoadIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	bumped := strings.Replace(
		string(readEditorStateFixture(t, currentStateFixture)),
		`"contentDefaultsVersion": 1`,
		`"contentDefaultsVersion": 2`,
		1)
	statePath := filepath.Join(t.TempDir(), "tooNewDefaults.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(bumped), 0o644))
	target := mappers.NewEditorStateMapper().NewDefaultEntity()

	// Act
	_, err := newEditorStateMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedContentDefaultsVersionError))
}

// Files the user already has on disk span every older format. If any of them
// stops loading, the migration is wrong, not the file. output/ is gitignored,
// so this is an opt-in local check rather than a CI gate.
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
		_, loadErr := migrator.Load(path, &target)
		assert.NoError(t, loadErr, "%s no longer loads", filepath.Base(path))
	}
}

func TestWhenBothFlatStateFixturesAreParsed_TheyShareEveryKeyButTheSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	unversionedKeys := slices.Sorted(maps.Keys(parseEditorStateFixtureKeys(t, unversionedFixture)))
	flatKeys := parseEditorStateFixtureKeys(t, previousStateFixture)

	// Act
	delete(flatKeys, schemaVersionKey)

	// Assert
	assert.Equal(t, unversionedKeys, slices.Sorted(maps.Keys(flatKeys)))
}

func TestWhenTheAllFieldsStateIsWritten_ItMatchesTheCurrentFixture(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := map[string]any{}
	require.NoError(t, json.Unmarshal(readEditorStateFixture(t, currentStateFixture), &expected))

	// Act
	//nolint:musttag // data.Vec2 is deliberately tag-free: its field names X and Y are the wire form.
	written, err := json.Marshal(test_helpers.NewAllFieldsEditorStateEntity())

	// Assert
	require.NoError(t, err)
	actual := map[string]any{}
	require.NoError(t, json.Unmarshal(written, &actual))
	assert.Equal(t, expected, actual)
}

// newStampFreeAllFieldsEntity is the all-fields state as it comes back from a
// pre-v2 file: identical but for the generator stamps, which those files had
// nowhere to put.
func newStampFreeAllFieldsEntity() editor_state.EditorState {
	expected := test_helpers.NewAllFieldsEditorStateEntity()
	expected.ManualEdits.Zones[0].GeneratorPosition = nil
	expected.ManualEdits.Zones[0].GeneratorRing = nil

	return expected
}

// decodeEditorStateFixture reads a fixture the way the app does: seeded with the
// defaults and routed through the migrator, which is the only path that can
// read an older file at all.
func decodeEditorStateFixture(t *testing.T, name string) (editor_state.EditorState, error) {
	t.Helper()

	target := mappers.NewEditorStateMapper().NewDefaultEntity()
	_, err := newEditorStateMigrator().
		Load(filepath.Join("..", "test_helpers", "testdata", name), &target)

	return target, err
}

func newEditorStateMigrator() editor_state_migrator.IEditorStateMigrator {
	return editor_state_migrator.NewEditorStateMigrator(
		repositories.NewEditorStateRepository(),
		repositories.NewV2EditorStateRepository(),
		repositories.NewLegacyEditorStateRepository())
}

// readEditorStateFixture returns a committed .gen.json fixture. `_v0_` freezes
// the flat format before versioning, `_v1_` is the last one with array
// positions, `_v2_` the last flat one, and `_v3_` is what the current writer
// emits. Regrouping reorders keys, so assertions compare parsed values.
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
