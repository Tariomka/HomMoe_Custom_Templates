package editorStateEntity_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenAKeyIsAbsentFromTheFile_TheReceiverKeepsItsSeededValue(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{TemplateName: "Seeded"}

	// Act
	err := json.Unmarshal([]byte(`{"playerCount":5}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", entity.TemplateName)
}

func TestWhenAKeyIsPresentInTheFile_ItReplacesTheSeededValue(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{TemplateName: "Seeded"}

	// Act
	err := json.Unmarshal([]byte(`{"templateName":"From File"}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "From File", entity.TemplateName)
}

func TestWhenTheFileCarriesNoSchemaVersion_ItIsMigratedToTheCurrentVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{}

	// Act
	err := json.Unmarshal([]byte(`{"templateName":"Legacy"}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, entity.SchemaVersion)
}

func TestWhenTheFileCarriesAKnownSchemaVersion_ItIsKeptVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{}

	// Act
	err := json.Unmarshal([]byte(`{"schemaVersion":7}`), &entity)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 7, entity.SchemaVersion)
}

func TestWhenTheFileIsNotJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.EditorStateEntity{}

	// Act
	err := json.Unmarshal([]byte("{not json"), &entity)

	// Assert
	assert.Error(t, err)
}
