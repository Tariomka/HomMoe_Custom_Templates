package editorState_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"maps"
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheStateIsWritten_IdentityMapAndSchemaFieldsSitAtTheRoot(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{
		ManualEdits: editor_state.ManualEditSettings{Zones: []editor_state.ManualZoneEdit{{Name: "Spawn-A"}}},
	}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.ElementsMatch(t,
		[]string{
			"schemaVersion", "generatorVersion", "templateName", "gameMode", "mapSize", "experimentalMapSizes",
			"players", "neutralZones", "castles", "generation", "gameRules", "content", "manualEdits",
		},
		slices.Collect(maps.Keys(keys)))
}

func TestWhenThereAreNoManualEdits_TheManualEditsSectionIsNotWritten(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, keys, "manualEdits")
}

func TestWhenTheStateIsWritten_SettingsSitInsideTheirSection(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{Players: editor_state.PlayerSettings{PlayerCount: 5}}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, `5`, string(sectionKeys(t, keys, "players")["playerCount"]))
}

func TestWhenContentRowsAreNil_TheKeyIsNotWritten(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{Content: editor_state.ContentSettings{PlayerZoneContentRows: nil}}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.NotContains(t, sectionKeys(t, keys, "content"), "playerZoneContentRows")
}

// A cleared list must survive as [], otherwise the load path would mistake it
// for omitted defaults.
func TestWhenContentRowsAreEmpty_TheKeyIsWrittenAsAnEmptyList(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{
		Content: editor_state.ContentSettings{PlayerZoneContentRows: []editor_state.ZoneContentRow{}},
	}

	// Act
	keys := marshalToKeys(t, entity)

	// Assert
	assert.JSONEq(t, `[]`, string(sectionKeys(t, keys, "content")["playerZoneContentRows"]))
}

// The load path seeds the entity with the defaults before decoding, so this
// merge is what lets a key the file omits keep its default instead of
// collapsing to a zero value.
func TestWhenAKeyIsAbsentFromTheFile_TheSeededValueSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorState{
		TemplateName: "Seeded",
		Players:      editor_state.PlayerSettings{HeroCountMax: 8},
	}

	// Act
	err := json.Unmarshal([]byte(`{"players":{"playerCount":5}}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t,
		[]any{"Seeded", 8},
		[]any{entity.TemplateName, entity.Players.HeroCountMax})
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

func sectionKeys(t *testing.T, keys map[string]jsontext.Value, section string) map[string]jsontext.Value {
	t.Helper()

	nested := map[string]jsontext.Value{}
	require.NoError(t, json.Unmarshal(keys[section], &nested))

	return nested
}
