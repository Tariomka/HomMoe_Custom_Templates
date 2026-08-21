package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// persistedEditorStateFieldCount is the number of top-level keys the editor
// state writes to a .gen.json file.
const persistedEditorStateFieldCount = 72

func TestWhenTheFrozenStateFixtureIsLoaded_EveryPersistedFieldKeepsItsValue(t *testing.T) {
	t.Parallel()
	// Arrange
	data := readFrozenEditorStateFixture(t)
	var loaded editor_state_dto.EditorStateDto

	// Act
	err := json.Unmarshal(data, &loaded)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, test_helpers.NewAllFieldsEditorStateDto(), loaded)
}

func TestWhenTheFrozenStateFixtureIsParsed_ItCarriesEveryPersistedKey(t *testing.T) {
	t.Parallel()
	// Arrange
	data := readFrozenEditorStateFixture(t)
	keys := map[string]json.RawMessage{}

	// Act
	err := json.Unmarshal(data, &keys)

	// Assert
	require.NoError(t, err)
	assert.Len(t, keys, persistedEditorStateFieldCount)
}

func TestWhenTheAllFieldsStateIsWritten_ItMatchesTheFrozenFixture(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := map[string]any{}
	require.NoError(t, json.Unmarshal(readFrozenEditorStateFixture(t), &expected))

	// Act
	written, err := json.Marshal(test_helpers.NewAllFieldsEditorStateDto())

	// Assert
	require.NoError(t, err)
	actual := map[string]any{}
	require.NoError(t, json.Unmarshal(written, &actual))
	assert.Equal(t, expected, actual)
}

// readFrozenEditorStateFixture returns the committed .gen.json fixture that
// freezes the wire format produced by the editor before the editor state was
// split into entity groups. Regrouping the fields reorders the keys, so every
// assertion above compares parsed values and never raw bytes.
func readFrozenEditorStateFixture(t *testing.T) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "test_helpers", "testdata", "editorState_v0_flat.gen.json"))
	require.NoError(t, err)
	return data
}
