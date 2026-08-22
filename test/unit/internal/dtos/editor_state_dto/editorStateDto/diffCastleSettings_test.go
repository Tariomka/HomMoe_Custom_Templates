package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenHubCastleCountChanges_FlagsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	current := previous
	current.HubZoneCastles = 3

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_model.CastleSettingChanges{Hub: true}, changes)
}

func TestWhenNoCastleOptionChanges_ReportsNoChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	current := previous

	// Act
	changes := previous.DiffCastleSettings(&current)

	// Assert
	assert.Equal(t, editor_state_model.CastleSettingChanges{}, changes)
}
