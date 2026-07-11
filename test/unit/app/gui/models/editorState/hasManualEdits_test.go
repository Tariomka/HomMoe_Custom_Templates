package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoManualEditsWereSet_ReportsNoManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.False(t, hasEdits)
}

func TestWhenManualZonesWereStored_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SetManualEdits([]entities.Zone{{Name: "Zone A"}}, nil)

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
