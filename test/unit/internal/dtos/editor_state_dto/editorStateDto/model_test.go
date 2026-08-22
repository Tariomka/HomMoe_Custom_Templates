package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheCarriedModelIsRequested_ReturnsTheStateTheDtoWasBuiltFrom(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()
	dto := editor_state_dto.NewEditorStateDto(state)

	// Act
	carried := dto.Model()

	// Assert
	assert.Equal(t, state, *carried)
}

func TestWhenTheCarriedModelIsMutated_TheDtoSeesTheChange(t *testing.T) {
	t.Parallel()
	// Arrange
	dto := editor_state_dto.NewDefaultEditorStateDto()
	carried := dto.Model()
	require.NotEqual(t, "Renamed", dto.TemplateName)

	// Act
	carried.TemplateName = "Renamed"

	// Assert
	assert.Equal(t, "Renamed", dto.TemplateName)
}
