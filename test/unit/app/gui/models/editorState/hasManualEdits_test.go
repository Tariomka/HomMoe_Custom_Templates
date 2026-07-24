package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"testing"
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
	state.SetManualEdits([]entities.Zone{{Name: "Zone A"}}, nil)

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
