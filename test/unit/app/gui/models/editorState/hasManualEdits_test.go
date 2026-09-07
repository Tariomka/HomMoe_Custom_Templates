package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoManualEditsWereSet_ReportsNoManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.False(t, hasEdits)
}

func TestWhenManualZonesWereStored_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits([]template_model.Zone{{Name: "Zone A"}}, nil)

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
